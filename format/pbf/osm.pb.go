// Code generated by protoc-gen-go. DO NOT EDIT.
// source: format/pbf/osm.proto

package pbf

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Relation_MemberType int32

const (
	Relation_NODE     Relation_MemberType = 0
	Relation_WAY      Relation_MemberType = 1
	Relation_RELATION Relation_MemberType = 2
)

var Relation_MemberType_name = map[int32]string{
	0: "NODE",
	1: "WAY",
	2: "RELATION",
}

var Relation_MemberType_value = map[string]int32{
	"NODE":     0,
	"WAY":      1,
	"RELATION": 2,
}

func (x Relation_MemberType) String() string {
	return proto.EnumName(Relation_MemberType_name, int32(x))
}

func (Relation_MemberType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_bbf78f6194d7fe88, []int{4, 0}
}

type BlobHeader struct {
	Type                 string   `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Indexdata            []byte   `protobuf:"bytes,2,opt,name=indexdata,proto3" json:"indexdata,omitempty"`
	Datasize             int32    `protobuf:"varint,3,opt,name=datasize,proto3" json:"datasize,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BlobHeader) Reset()         { *m = BlobHeader{} }
func (m *BlobHeader) String() string { return proto.CompactTextString(m) }
func (*BlobHeader) ProtoMessage()    {}
func (*BlobHeader) Descriptor() ([]byte, []int) {
	return fileDescriptor_bbf78f6194d7fe88, []int{0}
}

func (m *BlobHeader) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BlobHeader.Unmarshal(m, b)
}
func (m *BlobHeader) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BlobHeader.Marshal(b, m, deterministic)
}
func (m *BlobHeader) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BlobHeader.Merge(m, src)
}
func (m *BlobHeader) XXX_Size() int {
	return xxx_messageInfo_BlobHeader.Size(m)
}
func (m *BlobHeader) XXX_DiscardUnknown() {
	xxx_messageInfo_BlobHeader.DiscardUnknown(m)
}

var xxx_messageInfo_BlobHeader proto.InternalMessageInfo

func (m *BlobHeader) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *BlobHeader) GetIndexdata() []byte {
	if m != nil {
		return m.Indexdata
	}
	return nil
}

func (m *BlobHeader) GetDatasize() int32 {
	if m != nil {
		return m.Datasize
	}
	return 0
}

type Blob struct {
	Raw                  []byte   `protobuf:"bytes,1,opt,name=raw,proto3" json:"raw,omitempty"`
	RawSize              int32    `protobuf:"varint,2,opt,name=raw_size,json=rawSize,proto3" json:"raw_size,omitempty"`
	ZlibData             []byte   `protobuf:"bytes,3,opt,name=zlib_data,json=zlibData,proto3" json:"zlib_data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Blob) Reset()         { *m = Blob{} }
func (m *Blob) String() string { return proto.CompactTextString(m) }
func (*Blob) ProtoMessage()    {}
func (*Blob) Descriptor() ([]byte, []int) {
	return fileDescriptor_bbf78f6194d7fe88, []int{1}
}

func (m *Blob) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Blob.Unmarshal(m, b)
}
func (m *Blob) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Blob.Marshal(b, m, deterministic)
}
func (m *Blob) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Blob.Merge(m, src)
}
func (m *Blob) XXX_Size() int {
	return xxx_messageInfo_Blob.Size(m)
}
func (m *Blob) XXX_DiscardUnknown() {
	xxx_messageInfo_Blob.DiscardUnknown(m)
}

var xxx_messageInfo_Blob proto.InternalMessageInfo

func (m *Blob) GetRaw() []byte {
	if m != nil {
		return m.Raw
	}
	return nil
}

func (m *Blob) GetRawSize() int32 {
	if m != nil {
		return m.RawSize
	}
	return 0
}

func (m *Blob) GetZlibData() []byte {
	if m != nil {
		return m.ZlibData
	}
	return nil
}

type Node struct {
	Id                   int64    `protobuf:"zigzag64,1,opt,name=id,proto3" json:"id,omitempty"`
	Keys                 []uint32 `protobuf:"varint,2,rep,packed,name=keys,proto3" json:"keys,omitempty"`
	Vals                 []uint32 `protobuf:"varint,3,rep,packed,name=vals,proto3" json:"vals,omitempty"`
	Info                 *Info    `protobuf:"bytes,4,opt,name=info,proto3" json:"info,omitempty"`
	Lat                  int64    `protobuf:"zigzag64,8,opt,name=lat,proto3" json:"lat,omitempty"`
	Lon                  int64    `protobuf:"zigzag64,9,opt,name=lon,proto3" json:"lon,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Node) Reset()         { *m = Node{} }
func (m *Node) String() string { return proto.CompactTextString(m) }
func (*Node) ProtoMessage()    {}
func (*Node) Descriptor() ([]byte, []int) {
	return fileDescriptor_bbf78f6194d7fe88, []int{2}
}

func (m *Node) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Node.Unmarshal(m, b)
}
func (m *Node) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Node.Marshal(b, m, deterministic)
}
func (m *Node) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Node.Merge(m, src)
}
func (m *Node) XXX_Size() int {
	return xxx_messageInfo_Node.Size(m)
}
func (m *Node) XXX_DiscardUnknown() {
	xxx_messageInfo_Node.DiscardUnknown(m)
}

var xxx_messageInfo_Node proto.InternalMessageInfo

func (m *Node) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Node) GetKeys() []uint32 {
	if m != nil {
		return m.Keys
	}
	return nil
}

func (m *Node) GetVals() []uint32 {
	if m != nil {
		return m.Vals
	}
	return nil
}

func (m *Node) GetInfo() *Info {
	if m != nil {
		return m.Info
	}
	return nil
}

func (m *Node) GetLat() int64 {
	if m != nil {
		return m.Lat
	}
	return 0
}

func (m *Node) GetLon() int64 {
	if m != nil {
		return m.Lon
	}
	return 0
}

type Way struct {
	Id                   int64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Keys                 []uint32 `protobuf:"varint,2,rep,packed,name=keys,proto3" json:"keys,omitempty"`
	Vals                 []uint32 `protobuf:"varint,3,rep,packed,name=vals,proto3" json:"vals,omitempty"`
	Info                 *Info    `protobuf:"bytes,4,opt,name=info,proto3" json:"info,omitempty"`
	Refs                 []int64  `protobuf:"zigzag64,8,rep,packed,name=refs,proto3" json:"refs,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Way) Reset()         { *m = Way{} }
func (m *Way) String() string { return proto.CompactTextString(m) }
func (*Way) ProtoMessage()    {}
func (*Way) Descriptor() ([]byte, []int) {
	return fileDescriptor_bbf78f6194d7fe88, []int{3}
}

func (m *Way) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Way.Unmarshal(m, b)
}
func (m *Way) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Way.Marshal(b, m, deterministic)
}
func (m *Way) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Way.Merge(m, src)
}
func (m *Way) XXX_Size() int {
	return xxx_messageInfo_Way.Size(m)
}
func (m *Way) XXX_DiscardUnknown() {
	xxx_messageInfo_Way.DiscardUnknown(m)
}

var xxx_messageInfo_Way proto.InternalMessageInfo

func (m *Way) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Way) GetKeys() []uint32 {
	if m != nil {
		return m.Keys
	}
	return nil
}

func (m *Way) GetVals() []uint32 {
	if m != nil {
		return m.Vals
	}
	return nil
}

func (m *Way) GetInfo() *Info {
	if m != nil {
		return m.Info
	}
	return nil
}

func (m *Way) GetRefs() []int64 {
	if m != nil {
		return m.Refs
	}
	return nil
}

type Relation struct {
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// Parallel arrays.
	Keys []uint32 `protobuf:"varint,2,rep,packed,name=keys,proto3" json:"keys,omitempty"`
	Vals []uint32 `protobuf:"varint,3,rep,packed,name=vals,proto3" json:"vals,omitempty"`
	Info *Info    `protobuf:"bytes,4,opt,name=info,proto3" json:"info,omitempty"`
	// Parallel arrays
	RolesSid             []int32               `protobuf:"varint,8,rep,packed,name=roles_sid,json=rolesSid,proto3" json:"roles_sid,omitempty"`
	Memids               []int64               `protobuf:"zigzag64,9,rep,packed,name=memids,proto3" json:"memids,omitempty"`
	Types                []Relation_MemberType `protobuf:"varint,10,rep,packed,name=types,proto3,enum=pbf.Relation_MemberType" json:"types,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *Relation) Reset()         { *m = Relation{} }
func (m *Relation) String() string { return proto.CompactTextString(m) }
func (*Relation) ProtoMessage()    {}
func (*Relation) Descriptor() ([]byte, []int) {
	return fileDescriptor_bbf78f6194d7fe88, []int{4}
}

func (m *Relation) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Relation.Unmarshal(m, b)
}
func (m *Relation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Relation.Marshal(b, m, deterministic)
}
func (m *Relation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Relation.Merge(m, src)
}
func (m *Relation) XXX_Size() int {
	return xxx_messageInfo_Relation.Size(m)
}
func (m *Relation) XXX_DiscardUnknown() {
	xxx_messageInfo_Relation.DiscardUnknown(m)
}

var xxx_messageInfo_Relation proto.InternalMessageInfo

func (m *Relation) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Relation) GetKeys() []uint32 {
	if m != nil {
		return m.Keys
	}
	return nil
}

func (m *Relation) GetVals() []uint32 {
	if m != nil {
		return m.Vals
	}
	return nil
}

func (m *Relation) GetInfo() *Info {
	if m != nil {
		return m.Info
	}
	return nil
}

func (m *Relation) GetRolesSid() []int32 {
	if m != nil {
		return m.RolesSid
	}
	return nil
}

func (m *Relation) GetMemids() []int64 {
	if m != nil {
		return m.Memids
	}
	return nil
}

func (m *Relation) GetTypes() []Relation_MemberType {
	if m != nil {
		return m.Types
	}
	return nil
}

type HeaderBBox struct {
	Left                 int64    `protobuf:"zigzag64,1,opt,name=left,proto3" json:"left,omitempty"`
	Right                int64    `protobuf:"zigzag64,2,opt,name=right,proto3" json:"right,omitempty"`
	Top                  int64    `protobuf:"zigzag64,3,opt,name=top,proto3" json:"top,omitempty"`
	Bottom               int64    `protobuf:"zigzag64,4,opt,name=bottom,proto3" json:"bottom,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HeaderBBox) Reset()         { *m = HeaderBBox{} }
func (m *HeaderBBox) String() string { return proto.CompactTextString(m) }
func (*HeaderBBox) ProtoMessage()    {}
func (*HeaderBBox) Descriptor() ([]byte, []int) {
	return fileDescriptor_bbf78f6194d7fe88, []int{5}
}

func (m *HeaderBBox) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HeaderBBox.Unmarshal(m, b)
}
func (m *HeaderBBox) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HeaderBBox.Marshal(b, m, deterministic)
}
func (m *HeaderBBox) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HeaderBBox.Merge(m, src)
}
func (m *HeaderBBox) XXX_Size() int {
	return xxx_messageInfo_HeaderBBox.Size(m)
}
func (m *HeaderBBox) XXX_DiscardUnknown() {
	xxx_messageInfo_HeaderBBox.DiscardUnknown(m)
}

var xxx_messageInfo_HeaderBBox proto.InternalMessageInfo

func (m *HeaderBBox) GetLeft() int64 {
	if m != nil {
		return m.Left
	}
	return 0
}

func (m *HeaderBBox) GetRight() int64 {
	if m != nil {
		return m.Right
	}
	return 0
}

func (m *HeaderBBox) GetTop() int64 {
	if m != nil {
		return m.Top
	}
	return 0
}

func (m *HeaderBBox) GetBottom() int64 {
	if m != nil {
		return m.Bottom
	}
	return 0
}

type HeaderBlock struct {
	Bbox *HeaderBBox `protobuf:"bytes,1,opt,name=bbox,proto3" json:"bbox,omitempty"`
	// Additional tags to aid in parsing this dataset
	RequiredFeatures []string `protobuf:"bytes,4,rep,name=required_features,json=requiredFeatures,proto3" json:"required_features,omitempty"`
	OptionalFeatures []string `protobuf:"bytes,5,rep,name=optional_features,json=optionalFeatures,proto3" json:"optional_features,omitempty"`
	Writingprogram   string   `protobuf:"bytes,16,opt,name=writingprogram,proto3" json:"writingprogram,omitempty"`
	Source           string   `protobuf:"bytes,17,opt,name=source,proto3" json:"source,omitempty"`
	// replication timestamp, expressed in seconds since the epoch,
	// otherwise the same value as in the "timestamp=..." field
	// in the state.txt file used by Osmosis
	OsmosisReplicationTimestamp int64 `protobuf:"varint,32,opt,name=osmosis_replication_timestamp,json=osmosisReplicationTimestamp,proto3" json:"osmosis_replication_timestamp,omitempty"`
	// replication sequence number (sequenceNumber in state.txt)
	OsmosisReplicationSequenceNumber int64 `protobuf:"varint,33,opt,name=osmosis_replication_sequence_number,json=osmosisReplicationSequenceNumber,proto3" json:"osmosis_replication_sequence_number,omitempty"`
	// replication base URL (from Osmosis' configuration.txt file)
	OsmosisReplicationBaseUrl string   `protobuf:"bytes,34,opt,name=osmosis_replication_base_url,json=osmosisReplicationBaseUrl,proto3" json:"osmosis_replication_base_url,omitempty"`
	XXX_NoUnkeyedLiteral      struct{} `json:"-"`
	XXX_unrecognized          []byte   `json:"-"`
	XXX_sizecache             int32    `json:"-"`
}

func (m *HeaderBlock) Reset()         { *m = HeaderBlock{} }
func (m *HeaderBlock) String() string { return proto.CompactTextString(m) }
func (*HeaderBlock) ProtoMessage()    {}
func (*HeaderBlock) Descriptor() ([]byte, []int) {
	return fileDescriptor_bbf78f6194d7fe88, []int{6}
}

func (m *HeaderBlock) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HeaderBlock.Unmarshal(m, b)
}
func (m *HeaderBlock) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HeaderBlock.Marshal(b, m, deterministic)
}
func (m *HeaderBlock) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HeaderBlock.Merge(m, src)
}
func (m *HeaderBlock) XXX_Size() int {
	return xxx_messageInfo_HeaderBlock.Size(m)
}
func (m *HeaderBlock) XXX_DiscardUnknown() {
	xxx_messageInfo_HeaderBlock.DiscardUnknown(m)
}

var xxx_messageInfo_HeaderBlock proto.InternalMessageInfo

func (m *HeaderBlock) GetBbox() *HeaderBBox {
	if m != nil {
		return m.Bbox
	}
	return nil
}

func (m *HeaderBlock) GetRequiredFeatures() []string {
	if m != nil {
		return m.RequiredFeatures
	}
	return nil
}

func (m *HeaderBlock) GetOptionalFeatures() []string {
	if m != nil {
		return m.OptionalFeatures
	}
	return nil
}

func (m *HeaderBlock) GetWritingprogram() string {
	if m != nil {
		return m.Writingprogram
	}
	return ""
}

func (m *HeaderBlock) GetSource() string {
	if m != nil {
		return m.Source
	}
	return ""
}

func (m *HeaderBlock) GetOsmosisReplicationTimestamp() int64 {
	if m != nil {
		return m.OsmosisReplicationTimestamp
	}
	return 0
}

func (m *HeaderBlock) GetOsmosisReplicationSequenceNumber() int64 {
	if m != nil {
		return m.OsmosisReplicationSequenceNumber
	}
	return 0
}

func (m *HeaderBlock) GetOsmosisReplicationBaseUrl() string {
	if m != nil {
		return m.OsmosisReplicationBaseUrl
	}
	return ""
}

type PrimitiveBlock struct {
	Stringtable    *StringTable      `protobuf:"bytes,1,opt,name=stringtable,proto3" json:"stringtable,omitempty"`
	Primitivegroup []*PrimitiveGroup `protobuf:"bytes,2,rep,name=primitivegroup,proto3" json:"primitivegroup,omitempty"`
	// Granularity, units of nanodegrees, used to store coordinates in this block
	Granularity int32 `protobuf:"varint,17,opt,name=granularity,proto3" json:"granularity,omitempty"`
	// Offset value between the output coordinates coordinates and the granularity grid, in units of nanodegrees.
	LatOffset int64 `protobuf:"varint,19,opt,name=lat_offset,json=latOffset,proto3" json:"lat_offset,omitempty"`
	LonOffset int64 `protobuf:"varint,20,opt,name=lon_offset,json=lonOffset,proto3" json:"lon_offset,omitempty"`
	// Granularity of dates, normally represented in units of milliseconds since the 1970 epoch.
	DateGranularity      int32    `protobuf:"varint,18,opt,name=date_granularity,json=dateGranularity,proto3" json:"date_granularity,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PrimitiveBlock) Reset()         { *m = PrimitiveBlock{} }
func (m *PrimitiveBlock) String() string { return proto.CompactTextString(m) }
func (*PrimitiveBlock) ProtoMessage()    {}
func (*PrimitiveBlock) Descriptor() ([]byte, []int) {
	return fileDescriptor_bbf78f6194d7fe88, []int{7}
}

func (m *PrimitiveBlock) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PrimitiveBlock.Unmarshal(m, b)
}
func (m *PrimitiveBlock) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PrimitiveBlock.Marshal(b, m, deterministic)
}
func (m *PrimitiveBlock) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PrimitiveBlock.Merge(m, src)
}
func (m *PrimitiveBlock) XXX_Size() int {
	return xxx_messageInfo_PrimitiveBlock.Size(m)
}
func (m *PrimitiveBlock) XXX_DiscardUnknown() {
	xxx_messageInfo_PrimitiveBlock.DiscardUnknown(m)
}

var xxx_messageInfo_PrimitiveBlock proto.InternalMessageInfo

func (m *PrimitiveBlock) GetStringtable() *StringTable {
	if m != nil {
		return m.Stringtable
	}
	return nil
}

func (m *PrimitiveBlock) GetPrimitivegroup() []*PrimitiveGroup {
	if m != nil {
		return m.Primitivegroup
	}
	return nil
}

func (m *PrimitiveBlock) GetGranularity() int32 {
	if m != nil {
		return m.Granularity
	}
	return 0
}

func (m *PrimitiveBlock) GetLatOffset() int64 {
	if m != nil {
		return m.LatOffset
	}
	return 0
}

func (m *PrimitiveBlock) GetLonOffset() int64 {
	if m != nil {
		return m.LonOffset
	}
	return 0
}

func (m *PrimitiveBlock) GetDateGranularity() int32 {
	if m != nil {
		return m.DateGranularity
	}
	return 0
}

type ChangeSet struct {
	Id                   int64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ChangeSet) Reset()         { *m = ChangeSet{} }
func (m *ChangeSet) String() string { return proto.CompactTextString(m) }
func (*ChangeSet) ProtoMessage()    {}
func (*ChangeSet) Descriptor() ([]byte, []int) {
	return fileDescriptor_bbf78f6194d7fe88, []int{8}
}

func (m *ChangeSet) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChangeSet.Unmarshal(m, b)
}
func (m *ChangeSet) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChangeSet.Marshal(b, m, deterministic)
}
func (m *ChangeSet) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChangeSet.Merge(m, src)
}
func (m *ChangeSet) XXX_Size() int {
	return xxx_messageInfo_ChangeSet.Size(m)
}
func (m *ChangeSet) XXX_DiscardUnknown() {
	xxx_messageInfo_ChangeSet.DiscardUnknown(m)
}

var xxx_messageInfo_ChangeSet proto.InternalMessageInfo

func (m *ChangeSet) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

type PrimitiveGroup struct {
	Nodes                []*Node      `protobuf:"bytes,1,rep,name=nodes,proto3" json:"nodes,omitempty"`
	Dense                *DenseNodes  `protobuf:"bytes,2,opt,name=dense,proto3" json:"dense,omitempty"`
	Ways                 []*Way       `protobuf:"bytes,3,rep,name=ways,proto3" json:"ways,omitempty"`
	Relations            []*Relation  `protobuf:"bytes,4,rep,name=relations,proto3" json:"relations,omitempty"`
	Changesets           []*ChangeSet `protobuf:"bytes,5,rep,name=changesets,proto3" json:"changesets,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *PrimitiveGroup) Reset()         { *m = PrimitiveGroup{} }
func (m *PrimitiveGroup) String() string { return proto.CompactTextString(m) }
func (*PrimitiveGroup) ProtoMessage()    {}
func (*PrimitiveGroup) Descriptor() ([]byte, []int) {
	return fileDescriptor_bbf78f6194d7fe88, []int{9}
}

func (m *PrimitiveGroup) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PrimitiveGroup.Unmarshal(m, b)
}
func (m *PrimitiveGroup) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PrimitiveGroup.Marshal(b, m, deterministic)
}
func (m *PrimitiveGroup) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PrimitiveGroup.Merge(m, src)
}
func (m *PrimitiveGroup) XXX_Size() int {
	return xxx_messageInfo_PrimitiveGroup.Size(m)
}
func (m *PrimitiveGroup) XXX_DiscardUnknown() {
	xxx_messageInfo_PrimitiveGroup.DiscardUnknown(m)
}

var xxx_messageInfo_PrimitiveGroup proto.InternalMessageInfo

func (m *PrimitiveGroup) GetNodes() []*Node {
	if m != nil {
		return m.Nodes
	}
	return nil
}

func (m *PrimitiveGroup) GetDense() *DenseNodes {
	if m != nil {
		return m.Dense
	}
	return nil
}

func (m *PrimitiveGroup) GetWays() []*Way {
	if m != nil {
		return m.Ways
	}
	return nil
}

func (m *PrimitiveGroup) GetRelations() []*Relation {
	if m != nil {
		return m.Relations
	}
	return nil
}

func (m *PrimitiveGroup) GetChangesets() []*ChangeSet {
	if m != nil {
		return m.Changesets
	}
	return nil
}

type StringTable struct {
	S                    []string `protobuf:"bytes,1,rep,name=s,proto3" json:"s,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StringTable) Reset()         { *m = StringTable{} }
func (m *StringTable) String() string { return proto.CompactTextString(m) }
func (*StringTable) ProtoMessage()    {}
func (*StringTable) Descriptor() ([]byte, []int) {
	return fileDescriptor_bbf78f6194d7fe88, []int{10}
}

func (m *StringTable) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StringTable.Unmarshal(m, b)
}
func (m *StringTable) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StringTable.Marshal(b, m, deterministic)
}
func (m *StringTable) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StringTable.Merge(m, src)
}
func (m *StringTable) XXX_Size() int {
	return xxx_messageInfo_StringTable.Size(m)
}
func (m *StringTable) XXX_DiscardUnknown() {
	xxx_messageInfo_StringTable.DiscardUnknown(m)
}

var xxx_messageInfo_StringTable proto.InternalMessageInfo

func (m *StringTable) GetS() []string {
	if m != nil {
		return m.S
	}
	return nil
}

type Info struct {
	Version   int32 `protobuf:"varint,1,opt,name=version,proto3" json:"version,omitempty"`
	Timestamp int32 `protobuf:"varint,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Changeset int64 `protobuf:"varint,3,opt,name=changeset,proto3" json:"changeset,omitempty"`
	Uid       int32 `protobuf:"varint,4,opt,name=uid,proto3" json:"uid,omitempty"`
	UserSid   int32 `protobuf:"varint,5,opt,name=user_sid,json=userSid,proto3" json:"user_sid,omitempty"`
	// The visible flag is used to store history information. It indicates that
	// the current object version has been created by a delete operation on the
	// OSM API.
	// When a writer sets this flag, it MUST add a required_features tag with
	// value "HistoricalInformation" to the HeaderBlock.
	// If this flag is not available for some object it MUST be assumed to be
	// true if the file has the required_features tag "HistoricalInformation"
	// set.
	Visible              bool     `protobuf:"varint,6,opt,name=visible,proto3" json:"visible,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Info) Reset()         { *m = Info{} }
func (m *Info) String() string { return proto.CompactTextString(m) }
func (*Info) ProtoMessage()    {}
func (*Info) Descriptor() ([]byte, []int) {
	return fileDescriptor_bbf78f6194d7fe88, []int{11}
}

func (m *Info) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Info.Unmarshal(m, b)
}
func (m *Info) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Info.Marshal(b, m, deterministic)
}
func (m *Info) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Info.Merge(m, src)
}
func (m *Info) XXX_Size() int {
	return xxx_messageInfo_Info.Size(m)
}
func (m *Info) XXX_DiscardUnknown() {
	xxx_messageInfo_Info.DiscardUnknown(m)
}

var xxx_messageInfo_Info proto.InternalMessageInfo

func (m *Info) GetVersion() int32 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *Info) GetTimestamp() int32 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *Info) GetChangeset() int64 {
	if m != nil {
		return m.Changeset
	}
	return 0
}

func (m *Info) GetUid() int32 {
	if m != nil {
		return m.Uid
	}
	return 0
}

func (m *Info) GetUserSid() int32 {
	if m != nil {
		return m.UserSid
	}
	return 0
}

func (m *Info) GetVisible() bool {
	if m != nil {
		return m.Visible
	}
	return false
}

type DenseNodes struct {
	Id []int64 `protobuf:"zigzag64,1,rep,packed,name=id,proto3" json:"id,omitempty"`
	//repeated Info info = 4;
	Denseinfo *DenseInfo `protobuf:"bytes,5,opt,name=denseinfo,proto3" json:"denseinfo,omitempty"`
	Lat       []int64    `protobuf:"zigzag64,8,rep,packed,name=lat,proto3" json:"lat,omitempty"`
	Lon       []int64    `protobuf:"zigzag64,9,rep,packed,name=lon,proto3" json:"lon,omitempty"`
	// Special packing of keys and vals into one array. May be empty if all nodes in this block are tagless.
	KeysVals             []int32  `protobuf:"varint,10,rep,packed,name=keys_vals,json=keysVals,proto3" json:"keys_vals,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DenseNodes) Reset()         { *m = DenseNodes{} }
func (m *DenseNodes) String() string { return proto.CompactTextString(m) }
func (*DenseNodes) ProtoMessage()    {}
func (*DenseNodes) Descriptor() ([]byte, []int) {
	return fileDescriptor_bbf78f6194d7fe88, []int{12}
}

func (m *DenseNodes) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DenseNodes.Unmarshal(m, b)
}
func (m *DenseNodes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DenseNodes.Marshal(b, m, deterministic)
}
func (m *DenseNodes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DenseNodes.Merge(m, src)
}
func (m *DenseNodes) XXX_Size() int {
	return xxx_messageInfo_DenseNodes.Size(m)
}
func (m *DenseNodes) XXX_DiscardUnknown() {
	xxx_messageInfo_DenseNodes.DiscardUnknown(m)
}

var xxx_messageInfo_DenseNodes proto.InternalMessageInfo

func (m *DenseNodes) GetId() []int64 {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *DenseNodes) GetDenseinfo() *DenseInfo {
	if m != nil {
		return m.Denseinfo
	}
	return nil
}

func (m *DenseNodes) GetLat() []int64 {
	if m != nil {
		return m.Lat
	}
	return nil
}

func (m *DenseNodes) GetLon() []int64 {
	if m != nil {
		return m.Lon
	}
	return nil
}

func (m *DenseNodes) GetKeysVals() []int32 {
	if m != nil {
		return m.KeysVals
	}
	return nil
}

type DenseInfo struct {
	Version   []int32 `protobuf:"varint,1,rep,packed,name=version,proto3" json:"version,omitempty"`
	Timestamp []int64 `protobuf:"zigzag64,2,rep,packed,name=timestamp,proto3" json:"timestamp,omitempty"`
	Changeset []int64 `protobuf:"zigzag64,3,rep,packed,name=changeset,proto3" json:"changeset,omitempty"`
	Uid       []int32 `protobuf:"zigzag32,4,rep,packed,name=uid,proto3" json:"uid,omitempty"`
	UserSid   []int32 `protobuf:"zigzag32,5,rep,packed,name=user_sid,json=userSid,proto3" json:"user_sid,omitempty"`
	// The visible flag is used to store history information. It indicates that
	// the current object version has been created by a delete operation on the
	// OSM API.
	// When a writer sets this flag, it MUST add a required_features tag with
	// value "HistoricalInformation" to the HeaderBlock.
	// If this flag is not available for some object it MUST be assumed to be
	// true if the file has the required_features tag "HistoricalInformation"
	// set.
	Visible              []bool   `protobuf:"varint,6,rep,packed,name=visible,proto3" json:"visible,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DenseInfo) Reset()         { *m = DenseInfo{} }
func (m *DenseInfo) String() string { return proto.CompactTextString(m) }
func (*DenseInfo) ProtoMessage()    {}
func (*DenseInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_bbf78f6194d7fe88, []int{13}
}

func (m *DenseInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DenseInfo.Unmarshal(m, b)
}
func (m *DenseInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DenseInfo.Marshal(b, m, deterministic)
}
func (m *DenseInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DenseInfo.Merge(m, src)
}
func (m *DenseInfo) XXX_Size() int {
	return xxx_messageInfo_DenseInfo.Size(m)
}
func (m *DenseInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_DenseInfo.DiscardUnknown(m)
}

var xxx_messageInfo_DenseInfo proto.InternalMessageInfo

func (m *DenseInfo) GetVersion() []int32 {
	if m != nil {
		return m.Version
	}
	return nil
}

func (m *DenseInfo) GetTimestamp() []int64 {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func (m *DenseInfo) GetChangeset() []int64 {
	if m != nil {
		return m.Changeset
	}
	return nil
}

func (m *DenseInfo) GetUid() []int32 {
	if m != nil {
		return m.Uid
	}
	return nil
}

func (m *DenseInfo) GetUserSid() []int32 {
	if m != nil {
		return m.UserSid
	}
	return nil
}

func (m *DenseInfo) GetVisible() []bool {
	if m != nil {
		return m.Visible
	}
	return nil
}

func init() {
	proto.RegisterEnum("pbf.Relation_MemberType", Relation_MemberType_name, Relation_MemberType_value)
	proto.RegisterType((*BlobHeader)(nil), "pbf.BlobHeader")
	proto.RegisterType((*Blob)(nil), "pbf.Blob")
	proto.RegisterType((*Node)(nil), "pbf.Node")
	proto.RegisterType((*Way)(nil), "pbf.Way")
	proto.RegisterType((*Relation)(nil), "pbf.Relation")
	proto.RegisterType((*HeaderBBox)(nil), "pbf.HeaderBBox")
	proto.RegisterType((*HeaderBlock)(nil), "pbf.HeaderBlock")
	proto.RegisterType((*PrimitiveBlock)(nil), "pbf.PrimitiveBlock")
	proto.RegisterType((*ChangeSet)(nil), "pbf.ChangeSet")
	proto.RegisterType((*PrimitiveGroup)(nil), "pbf.PrimitiveGroup")
	proto.RegisterType((*StringTable)(nil), "pbf.StringTable")
	proto.RegisterType((*Info)(nil), "pbf.Info")
	proto.RegisterType((*DenseNodes)(nil), "pbf.DenseNodes")
	proto.RegisterType((*DenseInfo)(nil), "pbf.DenseInfo")
}

func init() {
	proto.RegisterFile("format/pbf/osm.proto", fileDescriptor_bbf78f6194d7fe88)
}

var fileDescriptor_bbf78f6194d7fe88 = []byte{
	// 1089 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x56, 0x51, 0x6f, 0xdb, 0x36,
	0x10, 0x9e, 0x2c, 0x29, 0x91, 0xce, 0xa9, 0xeb, 0xb0, 0x46, 0xa0, 0x36, 0x09, 0xaa, 0xa9, 0xd8,
	0xe0, 0xa1, 0x5b, 0x0a, 0x78, 0x8f, 0x7b, 0x18, 0xea, 0xa5, 0xeb, 0x0a, 0xac, 0x49, 0xc1, 0x64,
	0x0b, 0xb6, 0x17, 0x8f, 0x8a, 0x68, 0x97, 0xa8, 0x2c, 0xaa, 0x24, 0x9d, 0xc4, 0xc1, 0xfe, 0x40,
	0x81, 0xfd, 0x83, 0xbd, 0x0c, 0xd8, 0x1f, 0xd9, 0x7f, 0xd8, 0x1f, 0x2a, 0x78, 0x92, 0x2c, 0x27,
	0xed, 0x6b, 0x9f, 0xac, 0xfb, 0xee, 0xe3, 0xdd, 0xf1, 0x3b, 0xf2, 0x68, 0x18, 0x4c, 0xa5, 0x9a,
	0x33, 0xf3, 0xa4, 0x4c, 0xa7, 0x4f, 0xa4, 0x9e, 0x1f, 0x94, 0x4a, 0x1a, 0x49, 0xdc, 0x32, 0x9d,
	0x26, 0xbf, 0x03, 0x8c, 0x73, 0x99, 0xfe, 0xc4, 0x59, 0xc6, 0x15, 0x21, 0xe0, 0x99, 0x65, 0xc9,
	0x23, 0x27, 0x76, 0x86, 0x21, 0xc5, 0x6f, 0xb2, 0x07, 0xa1, 0x28, 0x32, 0x7e, 0x95, 0x31, 0xc3,
	0xa2, 0x4e, 0xec, 0x0c, 0xb7, 0x68, 0x0b, 0x90, 0x07, 0x10, 0xd8, 0x5f, 0x2d, 0xae, 0x79, 0xe4,
	0xc6, 0xce, 0xd0, 0xa7, 0x2b, 0x3b, 0x79, 0x05, 0x9e, 0x8d, 0x4d, 0xfa, 0xe0, 0x2a, 0x76, 0x89,
	0x41, 0xb7, 0xa8, 0xfd, 0x24, 0xf7, 0x21, 0x50, 0xec, 0x72, 0x82, 0xab, 0x3a, 0xb8, 0x6a, 0x53,
	0xb1, 0xcb, 0x13, 0x71, 0xcd, 0xc9, 0x2e, 0x84, 0xd7, 0xb9, 0x48, 0x27, 0x98, 0xce, 0xc5, 0x25,
	0x81, 0x05, 0x0e, 0x99, 0x61, 0xc9, 0x3b, 0x07, 0xbc, 0x23, 0x99, 0x71, 0xd2, 0x83, 0x8e, 0xc8,
	0x30, 0x22, 0xa1, 0x1d, 0x91, 0x91, 0x1d, 0xf0, 0xde, 0xf0, 0xa5, 0x8e, 0x3a, 0xb1, 0x3b, 0xbc,
	0x33, 0xee, 0xf4, 0x1d, 0x8a, 0xb6, 0xc5, 0x2f, 0x58, 0xae, 0x23, 0xb7, 0xc5, 0xad, 0x4d, 0xf6,
	0xc1, 0x13, 0xc5, 0x54, 0x46, 0x5e, 0xec, 0x0c, 0xbb, 0xa3, 0xf0, 0xa0, 0x4c, 0xa7, 0x07, 0x2f,
	0x8a, 0xa9, 0xa4, 0x08, 0xdb, 0x8a, 0x73, 0x66, 0xa2, 0x00, 0xe3, 0xdb, 0x4f, 0x44, 0x64, 0x11,
	0x85, 0x35, 0x22, 0x8b, 0xe4, 0x4f, 0x70, 0xcf, 0xd8, 0x72, 0xad, 0x12, 0xf7, 0x53, 0x54, 0xb2,
	0x03, 0x9e, 0xe2, 0x53, 0x1d, 0x05, 0xb1, 0x3b, 0x24, 0xd5, 0x32, 0x6b, 0x27, 0xef, 0x3a, 0x10,
	0x50, 0x9e, 0x33, 0x23, 0x64, 0xf1, 0xa9, 0x6b, 0x78, 0x08, 0xa1, 0x92, 0x39, 0xd7, 0x13, 0x2d,
	0x32, 0x2c, 0xc4, 0xc7, 0xb5, 0x01, 0x82, 0x27, 0x22, 0x23, 0x0f, 0x60, 0x63, 0xce, 0xe7, 0x22,
	0xd3, 0x51, 0xb8, 0x2a, 0xb3, 0x46, 0xc8, 0x08, 0x7c, 0x7b, 0x8c, 0x74, 0x04, 0xb1, 0x3b, 0xec,
	0x8d, 0x22, 0x0c, 0xde, 0x54, 0x7e, 0xf0, 0x92, 0xcf, 0x53, 0xae, 0x4e, 0x97, 0x25, 0xc7, 0x45,
	0x15, 0x35, 0xf9, 0x06, 0xa0, 0x75, 0x90, 0x00, 0xbc, 0xa3, 0xe3, 0xc3, 0x67, 0xfd, 0xcf, 0xc8,
	0x26, 0xb8, 0x67, 0x4f, 0x7f, 0xeb, 0x3b, 0x64, 0x0b, 0x02, 0xfa, 0xec, 0xe7, 0xa7, 0xa7, 0x2f,
	0x8e, 0x8f, 0xfa, 0x9d, 0xe4, 0x0f, 0x80, 0xea, 0xfc, 0x8e, 0xc7, 0xf2, 0xca, 0x9e, 0xe1, 0x9c,
	0x4f, 0x4d, 0x7d, 0x38, 0xf0, 0x9b, 0x0c, 0xc0, 0x57, 0x62, 0xf6, 0xda, 0xe0, 0x61, 0x23, 0xb4,
	0x32, 0x6c, 0x4f, 0x8d, 0x2c, 0xf1, 0x90, 0x11, 0x6a, 0x3f, 0xc9, 0x0e, 0x6c, 0xa4, 0xd2, 0x18,
	0x39, 0x47, 0x29, 0x08, 0xad, 0xad, 0xe4, 0x1f, 0x17, 0xba, 0x75, 0x8a, 0x5c, 0x9e, 0xbf, 0x21,
	0x8f, 0xc0, 0x4b, 0x53, 0x79, 0x85, 0x39, 0xba, 0xa3, 0xbb, 0xb8, 0xa7, 0xb6, 0x04, 0x8a, 0x4e,
	0xf2, 0x18, 0xb6, 0x15, 0x7f, 0xbb, 0x10, 0x8a, 0x67, 0x93, 0x29, 0x67, 0x66, 0xa1, 0xb8, 0x8e,
	0xbc, 0xd8, 0x1d, 0x86, 0xb4, 0xdf, 0x38, 0x7e, 0xac, 0x71, 0x4b, 0x96, 0xa5, 0x95, 0x84, 0xe5,
	0x2d, 0xd9, 0xaf, 0xc8, 0x8d, 0x63, 0x45, 0xfe, 0x12, 0x7a, 0x97, 0x4a, 0x18, 0x51, 0xcc, 0x4a,
	0x25, 0x67, 0x8a, 0xcd, 0xa3, 0x3e, 0x5e, 0xd8, 0x5b, 0xa8, 0xdd, 0x8e, 0x96, 0x0b, 0x75, 0xce,
	0xa3, 0x6d, 0xf4, 0xd7, 0x16, 0x19, 0xc3, 0xbe, 0xd4, 0x73, 0xa9, 0x85, 0x9e, 0x28, 0x5e, 0xe6,
	0xe2, 0x1c, 0x9b, 0x31, 0x31, 0x62, 0xce, 0xb5, 0x61, 0xf3, 0x32, 0x8a, 0xf1, 0x28, 0xed, 0xd6,
	0x24, 0xda, 0x72, 0x4e, 0x1b, 0x0a, 0x79, 0x09, 0x8f, 0x3e, 0x16, 0x43, 0xf3, 0xb7, 0x0b, 0x5e,
	0x9c, 0xf3, 0x49, 0xb1, 0xb0, 0x0d, 0x8c, 0x3e, 0xc7, 0x48, 0xf1, 0x87, 0x91, 0x4e, 0x6a, 0xe2,
	0x11, 0xf2, 0xc8, 0xf7, 0xb0, 0xf7, 0xb1, 0x70, 0x29, 0xd3, 0x7c, 0xb2, 0x50, 0x79, 0x94, 0xe0,
	0x06, 0xee, 0x7f, 0x18, 0x67, 0xcc, 0x34, 0xff, 0x45, 0xe5, 0xc9, 0x5f, 0x1d, 0xe8, 0xbd, 0x52,
	0x62, 0x2e, 0x8c, 0xb8, 0xe0, 0x55, 0x97, 0x46, 0xd0, 0xd5, 0x46, 0x89, 0x62, 0x66, 0x58, 0x9a,
	0xf3, 0xba, 0x59, 0x7d, 0x6c, 0xd6, 0x09, 0xe2, 0xa7, 0x16, 0xa7, 0xeb, 0x24, 0xf2, 0x1d, 0xf4,
	0xca, 0x26, 0xca, 0x4c, 0xc9, 0x45, 0x89, 0x97, 0xa8, 0x3b, 0xba, 0x87, 0xcb, 0x56, 0x09, 0x9e,
	0x5b, 0x17, 0xbd, 0x45, 0x25, 0x31, 0x74, 0x67, 0x8a, 0x15, 0x8b, 0x9c, 0x29, 0x61, 0x96, 0x28,
	0xba, 0x4f, 0xd7, 0x21, 0xb2, 0x0f, 0x90, 0x33, 0x33, 0x91, 0xd3, 0xa9, 0xe6, 0x26, 0xba, 0x87,
	0xe2, 0x84, 0x39, 0x33, 0xc7, 0x08, 0xa0, 0x5b, 0x16, 0x8d, 0x7b, 0x50, 0xbb, 0x65, 0x51, 0xbb,
	0xbf, 0x82, 0x7e, 0xc6, 0x0c, 0x9f, 0xac, 0x27, 0x21, 0x98, 0xe4, 0xae, 0xc5, 0x9f, 0xb7, 0x70,
	0xb2, 0x0b, 0xe1, 0x0f, 0xaf, 0x59, 0x31, 0xe3, 0x27, 0xdc, 0xdc, 0x9e, 0x0f, 0xc9, 0xff, 0xce,
	0x9a, 0x56, 0xb8, 0x15, 0xf2, 0x10, 0xfc, 0x42, 0x66, 0x5c, 0x47, 0x0e, 0x6e, 0xb7, 0x9a, 0x01,
	0x76, 0xd4, 0xd2, 0x0a, 0x27, 0x5f, 0x80, 0x9f, 0xf1, 0x42, 0x57, 0xf3, 0xba, 0x39, 0xf3, 0x87,
	0x16, 0xb1, 0x2c, 0x4d, 0x2b, 0x2f, 0xd9, 0x03, 0xef, 0x92, 0x2d, 0xab, 0x11, 0xd3, 0x1d, 0x05,
	0xc8, 0x3a, 0x63, 0x4b, 0x8a, 0x28, 0x79, 0x0c, 0xa1, 0xaa, 0xaf, 0x7e, 0x75, 0x15, 0xba, 0xa3,
	0x3b, 0x37, 0x06, 0x02, 0x6d, 0xfd, 0xe4, 0x00, 0xe0, 0x1c, 0xb7, 0xa0, 0xb9, 0xa9, 0xee, 0x42,
	0x77, 0xd4, 0x43, 0xf6, 0x6a, 0x67, 0x74, 0x8d, 0x91, 0xec, 0x42, 0x77, 0xad, 0xad, 0x64, 0x0b,
	0x9c, 0x6a, 0x37, 0x21, 0x75, 0x74, 0xf2, 0xaf, 0x03, 0x9e, 0x1d, 0x69, 0x24, 0x82, 0xcd, 0x0b,
	0xae, 0xb4, 0x90, 0x05, 0x0a, 0xe2, 0xd3, 0xc6, 0xb4, 0x0f, 0x5d, 0x7b, 0x03, 0xaa, 0x57, 0xa9,
	0x05, 0xac, 0x77, 0x95, 0x0b, 0x47, 0x86, 0x4b, 0x5b, 0xc0, 0x8e, 0x92, 0x85, 0xc8, 0x70, 0x6a,
	0xf8, 0xd4, 0x7e, 0xda, 0x27, 0x6e, 0xa1, 0xb9, 0xc2, 0x99, 0xe9, 0x57, 0x89, 0xac, 0x6d, 0xc7,
	0xa5, 0x2d, 0x41, 0x68, 0x61, 0xcf, 0xe4, 0x46, 0xec, 0x0c, 0x03, 0xda, 0x98, 0xc9, 0xdf, 0x0e,
	0x40, 0xab, 0x29, 0x21, 0x75, 0xdf, 0x9a, 0x99, 0x6a, 0x67, 0xfb, 0xd7, 0x10, 0xa2, 0xd2, 0x38,
	0xb0, 0x7d, 0xec, 0x45, 0xaf, 0xed, 0x05, 0x4e, 0xed, 0x96, 0x40, 0x06, 0xcd, 0x43, 0xd6, 0x84,
	0xc0, 0xc7, 0x6c, 0xd0, 0x3c, 0x66, 0x2d, 0x2a, 0x0b, 0x3b, 0xe6, 0xed, 0x2b, 0x31, 0xc1, 0x27,
	0x02, 0xda, 0x31, 0x6f, 0xc1, 0x5f, 0x59, 0xae, 0x93, 0xff, 0x1c, 0x08, 0x57, 0x59, 0xc8, 0xde,
	0xba, 0x90, 0x0d, 0x79, 0x25, 0x66, 0x7c, 0x53, 0xcc, 0x26, 0xd1, 0x9a, 0xa0, 0xf1, 0x4d, 0x41,
	0x57, 0x8c, 0x56, 0xd4, 0x41, 0x23, 0xaa, 0x3b, 0xdc, 0xae, 0xca, 0xb4, 0xc2, 0xee, 0xdf, 0x10,
	0xb6, 0x71, 0xad, 0xc4, 0xdd, 0x5b, 0x17, 0xd7, 0x1d, 0x06, 0x75, 0x59, 0x15, 0x94, 0x6e, 0xe0,
	0x5f, 0x9f, 0x6f, 0xdf, 0x07, 0x00, 0x00, 0xff, 0xff, 0x34, 0x4b, 0x69, 0x9b, 0x12, 0x09, 0x00,
	0x00,
}
