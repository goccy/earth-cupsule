package osm

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"io"
	"runtime"
	"sync"
	"time"

	"github.com/goccy/earth-cupsule/format"
	"github.com/goccy/earth-cupsule/format/pbf"
	"github.com/gogo/protobuf/proto"
	"golang.org/x/xerrors"
)

// The length of the BlobHeader should be less than 32 KiB (32*1024 bytes) and must be less than 64 KiB.
// The uncompressed length of a Blob should be less than 16 MiB (16*1024*1024 bytes) and must be less than 32 MiB.
// ( See https://wiki.openstreetmap.org/wiki/PBF_Format )
const (
	maxBlobHeaderSize = 64 * 1024
	maxBlobSize       = 32 * 1024 * 1024
	bufSize           = 1024 * 1024
)

type BoundingBox struct {
	Left   float64
	Right  float64
	Top    float64
	Bottom float64
}

type Header struct {
	BoundingBox                      *BoundingBox
	RequiredFeatures                 []string
	OptionalFeatures                 []string
	WritingProgram                   string
	Source                           string
	OsmosisReplicationTimestamp      time.Time
	OsmosisReplicationSequenceNumber int64
	OsmosisReplicationBaseUrl        string
}

type Block struct {
	header *pbf.BlobHeader
	blob   *pbf.Blob
}

type Info struct {
	Version   int32
	Uid       int32
	Timestamp time.Time
	Changeset int64
	User      string
	Visible   bool
}

type decodeResult struct {
	offset int64
	err    error
}

type PBFDecoder struct {
	fsize         int64
	storage       *Storage
	offset        int64
	r             io.Reader
	buf           *bytes.Buffer
	header        *Header
	nodeCallback  func(*format.Node) error
	onceStart     sync.Once
	onceClose     sync.Once
	blockIndex    int
	concurrentNum int
	blocks        []chan *Block
}

func NewPBFDecoder(reader io.Reader, storage *Storage, fsize int64) *PBFDecoder {
	return &PBFDecoder{
		fsize:         fsize,
		storage:       storage,
		offset:        storage.Pos(),
		r:             reader,
		buf:           bytes.NewBuffer(make([]byte, 0, bufSize)),
		concurrentNum: runtime.NumCPU(),
	}
}

func (d *PBFDecoder) SetNodeCallback(cb func(*format.Node) error) {
	d.nodeCallback = cb
}

func (d *PBFDecoder) IsDecoded() bool {
	return d.offset >= d.fsize
}

func (d *PBFDecoder) Close() {
	d.onceClose.Do(func() {})
}

func (d *PBFDecoder) Decode() (offset int64, e error) {
	d.onceStart.Do(func() {
		if d.offset == 0 {
			header, size, err := d.readOSMHeader()
			if err != nil {
				e = xerrors.Errorf("failed to read osm header: %w", err)
			}
			d.offset += size
			d.header = header
		}

		block := make(chan *Block)
		for i := 0; i < d.concurrentNum; i++ {
			go func() {
				for blk := range block {
					if err := d.decodeAsOSMData(blk); err != nil {
						e = xerrors.Errorf("failed to decode as osm header: %w", err)
						d.Close()
					}
				}
			}()
			d.blocks = append(d.blocks, block)
		}
	})
	blk, size, err := d.readBlock()
	if err != nil {
		d.Close()
		return 0, xerrors.Errorf("failed to read block: %w", err)
	}

	block := d.blocks[d.blockIndex]
	d.blockIndex = (d.blockIndex + 1) % d.concurrentNum
	block <- blk
	d.offset += size
	return d.offset, nil
}

func (d *PBFDecoder) readBlock() (*Block, int64, error) {
	readSize := int64(4) // headerSize length
	headerSize, err := d.readBlobHeaderSize()
	if err != nil {
		return nil, 0, xerrors.Errorf("failed to read blob header size: %w", err)
	}

	readSize += int64(headerSize)
	header, err := d.readBlobHeader(headerSize)
	if err != nil {
		return nil, 0, xerrors.Errorf("failed to read header: %w", err)
	}

	readSize += int64(header.GetDatasize())
	blob, err := d.readBlob(header)
	if err != nil {
		return nil, 0, xerrors.Errorf("failed to read blob: %w", err)
	}
	return &Block{header: header, blob: blob}, readSize, nil
}

func (d *PBFDecoder) readBlobHeaderSize() (uint32, error) {
	d.buf.Reset()
	if _, err := io.CopyN(d.buf, d.r, 4); err != nil {
		return 0, xerrors.Errorf("failed to read blob header size: %w", err)
	}

	size := binary.BigEndian.Uint32(d.buf.Bytes())
	if size >= maxBlobHeaderSize {
		return 0, xerrors.Errorf("invalid blob header size %d", size)
	}
	return size, nil
}

func (d *PBFDecoder) readBlobHeader(headerSize uint32) (*pbf.BlobHeader, error) {
	d.buf.Reset()
	if _, err := io.CopyN(d.buf, d.r, int64(headerSize)); err != nil {
		return nil, xerrors.Errorf("failed to read blob header: %w", err)
	}

	var header pbf.BlobHeader
	if err := proto.Unmarshal(d.buf.Bytes(), &header); err != nil {
		return nil, xerrors.Errorf("failed to unmarshal blob header: %w", err)
	}

	if header.GetDatasize() >= maxBlobSize {
		return nil, xerrors.Errorf("invalid blob size %d", header.GetDatasize())
	}
	return &header, nil
}

func (d *PBFDecoder) readBlob(blobHeader *pbf.BlobHeader) (*pbf.Blob, error) {
	d.buf.Reset()
	if _, err := io.CopyN(d.buf, d.r, int64(blobHeader.GetDatasize())); err != nil {
		return nil, xerrors.Errorf("failed to read blob: %w", err)
	}

	var blob pbf.Blob
	if err := proto.Unmarshal(d.buf.Bytes(), &blob); err != nil {
		return nil, xerrors.Errorf("failed to unmarshal blob: %w", err)
	}
	return &blob, nil
}

func (d *PBFDecoder) blobData(blob *pbf.Blob) ([]byte, error) {
	if blob.Raw != nil {
		return blob.GetRaw(), nil
	}
	if blob.ZlibData != nil {
		r, err := zlib.NewReader(bytes.NewReader(blob.GetZlibData()))
		if err != nil {
			return nil, xerrors.Errorf("failed to create zlib reader: %w", err)
		}
		buf := bytes.NewBuffer(make([]byte, 0, blob.GetRawSize()+bytes.MinRead))
		if _, err := buf.ReadFrom(r); err != nil {
			return nil, xerrors.Errorf("faield to read zlib data: %w", err)
		}
		if buf.Len() != int(blob.GetRawSize()) {
			return nil, xerrors.Errorf("raw blob data size %d but expected %d", buf.Len(), blob.GetRawSize())
		}
		return buf.Bytes(), nil
	}
	return nil, errors.New("invalid blob data")
}

func (d *PBFDecoder) readOSMHeader() (*Header, int64, error) {
	block, size, err := d.readBlock()
	if err != nil {
		return nil, 0, xerrors.Errorf("failed to read block: %w", err)
	}
	header, err := d.decodeAsOSMHeader(block)
	if err != nil {
		return nil, 0, xerrors.Errorf("failed to decode as osm header: %w", err)
	}
	return header, size, nil
}

func (d *PBFDecoder) decodeAsOSMHeader(block *Block) (*Header, error) {
	data, err := d.blobData(block.blob)
	if err != nil {
		return nil, xerrors.Errorf("failed to get blob data: %w", err)
	}

	var header pbf.HeaderBlock
	if err := proto.Unmarshal(data, &header); err != nil {
		return nil, xerrors.Errorf("failed to unmarshal osm header: %w", err)
	}
	return NewHeader(&header), nil
}

func NewHeader(block *pbf.HeaderBlock) *Header {
	header := &Header{
		RequiredFeatures:                 block.GetRequiredFeatures(),
		OptionalFeatures:                 block.GetOptionalFeatures(),
		WritingProgram:                   block.GetWritingprogram(),
		Source:                           block.GetSource(),
		OsmosisReplicationBaseUrl:        block.GetOsmosisReplicationBaseUrl(),
		OsmosisReplicationSequenceNumber: block.GetOsmosisReplicationSequenceNumber(),
	}

	// convert timestamp epoch seconds to golang time structure if it exists
	if block.OsmosisReplicationTimestamp != 0 {
		header.OsmosisReplicationTimestamp = time.Unix(block.OsmosisReplicationTimestamp, 0)
	}
	// read bounding box if it exists
	if block.Bbox != nil {
		// Units are always in nanodegree and do not obey granularity rules. See osmformat.proto
		header.BoundingBox = &BoundingBox{
			Left:   1e-9 * float64(block.Bbox.Left),
			Right:  1e-9 * float64(block.Bbox.Right),
			Bottom: 1e-9 * float64(block.Bbox.Bottom),
			Top:    1e-9 * float64(block.Bbox.Top),
		}
	}
	return header
}

func (d *PBFDecoder) decodeOSMData() (int64, error) {
	block, size, err := d.readBlock()
	if err != nil {
		return 0, xerrors.Errorf("failed to read block: %w", err)
	}
	if err := d.decodeAsOSMData(block); err != nil {
		return 0, xerrors.Errorf("failed to decode as osm header: %w", err)
	}
	return size, nil
}

func (d *PBFDecoder) decodeAsOSMData(block *Block) error {
	data, err := d.blobData(block.blob)
	if err != nil {
		return xerrors.Errorf("failed to get blob data: %w", err)
	}

	var pblock pbf.PrimitiveBlock
	if err := proto.Unmarshal(data, &pblock); err != nil {
		return xerrors.Errorf("failed to unmarshal osm header: %w", err)
	}
	if err := d.decodePrimitiveBlock(&pblock); err != nil {
		return xerrors.Errorf("failed to decode primitive block: %w", err)
	}
	return nil
}

func (d *PBFDecoder) decodePrimitiveBlock(block *pbf.PrimitiveBlock) error {
	for _, group := range block.GetPrimitivegroup() {
		if err := d.decodePrimitiveGroup(block, group); err != nil {
			return xerrors.Errorf("failed to decode primitive group: %w", err)
		}
	}
	return nil
}

func (d *PBFDecoder) decodePrimitiveGroup(block *pbf.PrimitiveBlock, group *pbf.PrimitiveGroup) error {
	if err := d.decodeNodes(block, group.GetNodes()); err != nil {
		return xerrors.Errorf("failed to decode nodes: %w", err)
	}
	if err := d.decodeDenseNodes(block, group.GetDense()); err != nil {
		return xerrors.Errorf("failed to decode dense nodes: %w", err)
	}
	if err := d.decodeWays(block, group.GetWays()); err != nil {
		return xerrors.Errorf("failed to decode ways: %w", err)
	}
	if err := d.decodeRelations(block, group.GetRelations()); err != nil {
		return xerrors.Errorf("failed to decode relations: %w", err)
	}
	return nil
}

func (d *PBFDecoder) decodeNodes(block *pbf.PrimitiveBlock, nodes []*pbf.Node) error {
	st := block.GetStringtable().GetS()
	granularity := int64(block.GetGranularity())
	dateGranularity := int64(block.GetDateGranularity())

	latOffset := block.GetLatOffset()
	lonOffset := block.GetLonOffset()

	for _, node := range nodes {
		id := node.GetId()
		lat := node.GetLat()
		lon := node.GetLon()

		latitude := 1e-9 * float64((latOffset + (granularity * lat)))
		longitude := 1e-9 * float64((lonOffset + (granularity * lon)))

		tags := d.extractTags(st, node.GetKeys(), node.GetVals())
		info := d.extractInfo(st, node.GetInfo(), dateGranularity)

		node := &format.Node{
			ID:          id,
			Lat:         latitude,
			Lon:         longitude,
			User:        info.User,
			UserID:      int64(info.Uid),
			Visible:     info.Visible,
			Version:     int64(info.Version),
			ChangesetID: info.Changeset,
			Timestamp:   info.Timestamp,
			Tags:        format.Tags(tags),
		}
		if d.nodeCallback != nil && node.Tags.HasInterestingTags() {
			if err := d.nodeCallback(node); err != nil {
				return xerrors.Errorf("failed to node callback: %w", err)
			}
		}
		if err := d.storage.AddNode(node); err != nil {
			return xerrors.Errorf("failed to add node: %w", err)
		}
	}
	return nil
}

func (d *PBFDecoder) decodeDenseNodes(block *pbf.PrimitiveBlock, nodes *pbf.DenseNodes) error {
	st := block.GetStringtable().GetS()
	granularity := int64(block.GetGranularity())
	latOffset := block.GetLatOffset()
	lonOffset := block.GetLonOffset()
	dateGranularity := int64(block.GetDateGranularity())
	ids := nodes.GetId()
	lats := nodes.GetLat()
	lons := nodes.GetLon()
	di := nodes.GetDenseinfo()

	tu := tagUnpacker{st, nodes.GetKeysVals(), 0}
	var (
		id, lat, lon int64
		state        denseInfoState
	)
	for index := range ids {
		id = ids[index] + id
		lat = lats[index] + lat
		lon = lons[index] + lon
		latitude := 1e-9 * float64((latOffset + (granularity * lat)))
		longitude := 1e-9 * float64((lonOffset + (granularity * lon)))
		tags := tu.next()
		info := d.extractDenseInfo(st, &state, di, index, dateGranularity)

		node := &format.Node{
			ID:          id,
			Lat:         latitude,
			Lon:         longitude,
			User:        info.User,
			UserID:      int64(info.Uid),
			Visible:     info.Visible,
			Version:     int64(info.Version),
			ChangesetID: info.Changeset,
			Timestamp:   info.Timestamp,
			Tags:        format.Tags(tags),
		}
		if d.nodeCallback != nil && node.Tags.HasInterestingTags() {
			if err := d.nodeCallback(node); err != nil {
				return xerrors.Errorf("failed to node callback: %w", err)
			}
		}
		if err := d.storage.AddNode(node); err != nil {
			return xerrors.Errorf("failed to add node: %w", err)
		}
	}
	return nil
}

func (d *PBFDecoder) decodeWays(block *pbf.PrimitiveBlock, ways []*pbf.Way) error {
	st := block.GetStringtable().GetS()
	dateGranularity := int64(block.GetDateGranularity())

	for _, way := range ways {
		id := way.GetId()
		tags := d.extractTags(st, way.GetKeys(), way.GetVals())
		refs := way.GetRefs()
		var nodeID int64
		nodeRefs := make(format.NodeRefs, len(refs))
		for index := range refs {
			nodeID = refs[index] + nodeID // delta encoding
			nodeRefs[index] = &format.NodeRef{Ref: nodeID}
		}
		info := d.extractInfo(st, way.GetInfo(), dateGranularity)

		if err := d.storage.AddWay(&format.Way{
			ID:          id,
			User:        info.User,
			UserID:      int64(info.Uid),
			Visible:     info.Visible,
			ChangesetID: info.Changeset,
			Timestamp:   info.Timestamp,
			NodeRefs:    nodeRefs,
			Tags:        format.Tags(tags),
		}); err != nil {
			return xerrors.Errorf("failed to add way: %w", err)
		}
		for _, ref := range nodeRefs {
			if err := d.storage.AddNodeInWay(ref.Ref); err != nil {
				return xerrors.Errorf("failed to add node in way: %w", err)
			}
		}
	}
	return nil
}

func (d *PBFDecoder) decodeRelations(block *pbf.PrimitiveBlock, relations []*pbf.Relation) error {
	st := block.GetStringtable().GetS()
	dateGranularity := int64(block.GetDateGranularity())

	for _, rel := range relations {
		id := rel.GetId()
		tags := d.extractTags(st, rel.GetKeys(), rel.GetVals())
		members := d.extractMembers(st, rel)
		info := d.extractInfo(st, rel.GetInfo(), dateGranularity)

		if err := d.storage.AddRelation(&format.Relation{
			ID:          id,
			User:        info.User,
			UserID:      int64(info.Uid),
			Visible:     info.Visible,
			Version:     int64(info.Version),
			ChangesetID: info.Changeset,
			Timestamp:   info.Timestamp,
			Tags:        format.Tags(tags),
			Members:     members,
		}); err != nil {
			return xerrors.Errorf("failed to add relation: %w", err)
		}
		for _, m := range members {
			if err := d.storage.AddWayInMember(m.Ref); err != nil {
				return xerrors.Errorf("failed to add node in member: %w", err)
			}
		}
	}
	return nil
}

func (d *PBFDecoder) extractTags(stringTable []string, keyIDs, valueIDs []uint32) map[string]string {
	tags := make(map[string]string, len(keyIDs))
	for index, keyID := range keyIDs {
		key := stringTable[keyID]
		val := stringTable[valueIDs[index]]
		tags[key] = val
	}
	return tags
}

type tagUnpacker struct {
	stringTable []string
	keysVals    []int32
	index       int
}

// Make tags map from stringtable and array of IDs (used in DenseNodes encoding).
func (tu *tagUnpacker) next() map[string]string {
	tags := make(map[string]string)
	for tu.index < len(tu.keysVals) {
		keyID := tu.keysVals[tu.index]
		tu.index++
		if keyID == 0 {
			break
		}

		valID := tu.keysVals[tu.index]
		tu.index++

		key := tu.stringTable[keyID]
		val := tu.stringTable[valID]
		tags[key] = val
	}
	return tags
}

func (d *PBFDecoder) extractInfo(stringTable []string, i *pbf.Info, dateGranularity int64) *Info {
	info := &Info{Visible: true}
	if i == nil {
		return info
	}
	info.Version = i.GetVersion()

	millisec := time.Duration(int64(i.GetTimestamp())*dateGranularity) * time.Millisecond
	info.Timestamp = time.Unix(0, millisec.Nanoseconds()).UTC()
	info.Changeset = i.GetChangeset()
	info.Uid = i.GetUid()
	info.User = stringTable[i.GetUserSid()]
	info.Visible = i.GetVisible()
	return info
}

type denseInfoState struct {
	timestamp int64
	changeset int64
	uid       int32
	userSid   int32
}

func (d *PBFDecoder) extractDenseInfo(stringTable []string, state *denseInfoState, di *pbf.DenseInfo, index int, dateGranularity int64) Info {
	info := Info{Visible: true}

	versions := di.GetVersion()
	if len(versions) > 0 {
		info.Version = versions[index]
	}

	timestamps := di.GetTimestamp()
	if len(timestamps) > 0 {
		state.timestamp = timestamps[index] + state.timestamp
		millisec := time.Duration(state.timestamp*dateGranularity) * time.Millisecond
		info.Timestamp = time.Unix(0, millisec.Nanoseconds()).UTC()
	}

	changesets := di.GetChangeset()
	if len(changesets) > 0 {
		state.changeset = changesets[index] + state.changeset
		info.Changeset = state.changeset
	}

	uids := di.GetUid()
	if len(uids) > 0 {
		state.uid = uids[index] + state.uid
		info.Uid = state.uid
	}

	usersids := di.GetUserSid()
	if len(usersids) > 0 {
		state.userSid = usersids[index] + state.userSid
		info.User = stringTable[state.userSid]
	}

	visibleArray := di.GetVisible()
	if len(visibleArray) > 0 {
		info.Visible = visibleArray[index]
	}

	return info
}

// Make relation members from stringtable and three parallel arrays of IDs.
func (d *PBFDecoder) extractMembers(stringTable []string, rel *pbf.Relation) []*format.Member {
	memIDs := rel.GetMemids()
	types := rel.GetTypes()
	roleIDs := rel.GetRolesSid()

	var memID int64
	members := make([]*format.Member, len(memIDs))
	for index := range memIDs {
		memID = memIDs[index] + memID // delta encoding

		var memType format.Type
		switch types[index] {
		case pbf.Relation_NODE:
			memType = format.TypeNode
		case pbf.Relation_WAY:
			memType = format.TypeWay
		case pbf.Relation_RELATION:
			memType = format.TypeRelation
		}

		role := stringTable[roleIDs[index]]
		members[index] = &format.Member{
			Ref:  memID,
			Type: memType,
			Role: role,
		}
	}
	return members
}
