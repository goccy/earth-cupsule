package osm

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/goccy/earth-cupsule/format"
	"github.com/mitchellh/ioprogress"
	"golang.org/x/xerrors"
)

type Decoder struct {
	storage      *Storage
	file         *os.File
	xmlDecoder   *xml.Decoder
	lastElement  format.OSM
	isNeededStop bool
	offset       int64
	savedPos     int64
	onceClose    sync.Once
	nodeCallback func(*format.Node) error
	wayCallback  func(*format.Way) error
	relCallback  func(*format.Relation) error
}

func NewDecoder(path string) (*Decoder, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, xerrors.Errorf("failed to open file: %w", err)
	}
	finfo, err := f.Stat()
	if err != nil {
		return nil, xerrors.Errorf("failed to get stat: %w", err)
	}
	s, err := NewStorage(path)
	if err != nil {
		return nil, xerrors.Errorf("failed to open storage for osm: %w", err)
	}
	pos := s.Pos()
	if _, err := f.Seek(pos, os.SEEK_SET); err != nil {
		return nil, xerrors.Errorf("failed to seek: %w", err)
	}
	var progress *ioprogress.Reader
	if pos < finfo.Size() {
		progress = &ioprogress.Reader{
			Reader: f,
			Size:   finfo.Size(),
			DrawFunc: func(progress, total int64) error {
				progress += pos
				text := ioprogress.DrawTextFormatBytes(progress, total)
				if progress >= total {
					text = "done"
				}
				maxLineNum := 80
				if maxLineNum > len(text) {
					text += strings.Repeat(" ", maxLineNum-len(text))
				}
				_, err := fmt.Fprintf(os.Stdout, "progress... %s\r", text)
				return err
			},
		}
	}
	dec := &Decoder{
		storage:    s,
		file:       f,
		offset:     pos,
		xmlDecoder: xml.NewDecoder(progress),
	}
	return dec, nil
}

func (d *Decoder) close() error {
	if d.savedPos > 0 {
		if err := d.finish(); err != nil {
			return xerrors.Errorf("failed to finish: %w", err)
		}
	}
	if err := d.file.Close(); err != nil {
		return xerrors.Errorf("failed to close file: %w", err)
	}
	if err := d.storage.Close(); err != nil {
		return xerrors.Errorf("failed to close storage: %w", err)
	}
	return nil
}

func (d *Decoder) Close() (e error) {
	d.onceClose.Do(func() {
		if err := d.close(); err != nil {
			e = xerrors.Errorf("failed to close: %w", err)
		}
	})
	return
}

func (d *Decoder) Stop() {
	d.isNeededStop = true
}

func (d *Decoder) prepare() error {
	finfo, err := d.file.Stat()
	if err != nil {
		return xerrors.Errorf("failed to get stat: %w", err)
	}
	if d.offset >= finfo.Size() {
		return nil
	}
	for {
		if d.isNeededStop {
			return xerrors.New("force stop")
		}
		if err := d.decode(); err != nil {
			if xerrors.Is(err, io.EOF) {
				break
			}
			return xerrors.Errorf("failed to decode: %w", err)
		}
	}
	return nil
}

func (d *Decoder) setNextElement(e format.OSM) error {
	if d.lastElement != nil {
		switch e := d.lastElement.(type) {
		case *format.Node:
			if err := d.storage.AddNode(e); err != nil {
				return xerrors.Errorf("failed to add node: %w", err)
			}
		case *format.Way:
			if len(e.NodeRefs) == 0 {
				return xerrors.New("invalid way")
			}
			if err := d.storage.AddWay(e); err != nil {
				return xerrors.Errorf("failed to add way: %w", err)
			}
			for _, ref := range e.NodeRefs {
				if err := d.storage.AddNodeInWay(ref.Ref); err != nil {
					return xerrors.Errorf("failed to add node in way: %w", err)
				}
			}
		case *format.Relation:
			if len(e.Members) == 0 {
				return xerrors.New("invalid relation")
			}
			if err := d.storage.AddRelation(e); err != nil {
				return xerrors.Errorf("failed to add relation: %w", err)
			}
			for _, m := range e.Members {
				if err := d.storage.AddWayInMember(m.Ref); err != nil {
					return xerrors.Errorf("failed to add node in member: %w", err)
				}
			}
		}
	}
	d.lastElement = e
	return nil
}

func (d *Decoder) decodeNodeElement(e *xml.StartElement) (*format.Node, error) {
	node := format.NewNode()
	for _, attr := range e.Attr {
		switch attr.Name.Local {
		case "id":
			i, _ := strconv.ParseInt(attr.Value, 10, 64)
			node.ID = i
		case "lat":
			f, _ := strconv.ParseFloat(attr.Value, 64)
			node.Lat = f
		case "lon":
			f, _ := strconv.ParseFloat(attr.Value, 64)
			node.Lon = f
		case "user":
			node.User = attr.Value
		case "uid":
			i, _ := strconv.ParseInt(attr.Value, 10, 64)
			node.UserID = i
		case "visible":
			b, _ := strconv.ParseBool(attr.Value)
			node.Visible = b
		case "version":
			i, _ := strconv.ParseInt(attr.Value, 10, 64)
			node.Version = i
		case "changeset":
			i, _ := strconv.ParseInt(attr.Value, 10, 64)
			node.ChangesetID = i
		case "timestamp":
			t, _ := time.Parse("2006-1-2T15:4:5.9Z", attr.Value)
			node.Timestamp = t
		}
	}
	if err := d.setNextElement(node); err != nil {
		return nil, xerrors.Errorf("failed to set next node: %w", err)
	}
	return node, nil
}

func (d *Decoder) decodeWayElement(e *xml.StartElement) (*format.Way, error) {
	way := format.NewWay()
	for _, attr := range e.Attr {
		switch attr.Name.Local {
		case "id":
			i, _ := strconv.ParseInt(attr.Value, 10, 64)
			way.ID = i
		case "user":
			way.User = attr.Value
		case "uid":
			i, _ := strconv.ParseInt(attr.Value, 10, 64)
			way.UserID = i
		case "visible":
			b, _ := strconv.ParseBool(attr.Value)
			way.Visible = b
		case "version":
			i, _ := strconv.ParseInt(attr.Value, 10, 64)
			way.Version = i
		case "changeset":
			i, _ := strconv.ParseInt(attr.Value, 10, 64)
			way.ChangesetID = i
		case "timestamp":
			t, _ := time.Parse("2006-1-2T15:4:5.9Z", attr.Value)
			way.Timestamp = t
		}
	}
	if err := d.setNextElement(way); err != nil {
		return nil, xerrors.Errorf("failed to set next way: %w", err)
	}
	return way, nil
}

func (d *Decoder) decodeRelationElement(e *xml.StartElement) (*format.Relation, error) {
	rel := format.NewRelation()
	for _, attr := range e.Attr {
		switch attr.Name.Local {
		case "id":
			i, _ := strconv.ParseInt(attr.Value, 10, 64)
			rel.ID = i
		case "user":
			rel.User = attr.Value
		case "uid":
			i, _ := strconv.ParseInt(attr.Value, 10, 64)
			rel.UserID = i
		case "visible":
			b, _ := strconv.ParseBool(attr.Value)
			rel.Visible = b
		case "version":
			i, _ := strconv.ParseInt(attr.Value, 10, 64)
			rel.Version = i
		case "changeset":
			i, _ := strconv.ParseInt(attr.Value, 10, 64)
			rel.ChangesetID = i
		case "timestamp":
			t, _ := time.Parse("2006-1-2T15:4:5.9Z", attr.Value)
			rel.Timestamp = t
		}
	}
	if err := d.setNextElement(rel); err != nil {
		return nil, xerrors.Errorf("failed to set next relation: %w", err)
	}
	return rel, nil
}

func (d *Decoder) decodeMember(e *xml.StartElement) *format.Member {
	m := &format.Member{}
	for _, attr := range e.Attr {
		switch attr.Name.Local {
		case "type":
			m.Type = format.Type(attr.Value)
		case "ref":
			i, _ := strconv.ParseInt(attr.Value, 10, 64)
			m.Ref = i
		case "role":
			m.Role = attr.Value
		case "version":
			i, _ := strconv.ParseInt(attr.Value, 10, 64)
			m.Version = i
		case "changeset":
			i, _ := strconv.ParseInt(attr.Value, 10, 64)
			m.ChangesetID = i
		case "lat":
			f, _ := strconv.ParseFloat(attr.Value, 64)
			m.Lat = f
		case "lon":
			f, _ := strconv.ParseFloat(attr.Value, 64)
			m.Lon = f
		case "orientation":
			i, _ := strconv.ParseInt(attr.Value, 10, 64)
			m.Orientation = format.Orientation(i)
		}
	}
	return m
}

func (d *Decoder) finish() error {
	if err := d.storage.Finish(); err != nil {
		return xerrors.Errorf("failed to finish: %w", err)
	}
	if err := d.storage.SetPos(d.savedPos); err != nil {
		return xerrors.Errorf("failed to set pos: %w", err)
	}
	return nil
}

func (d *Decoder) decode() error {
	offset := d.offset + d.xmlDecoder.InputOffset()
	token, err := d.xmlDecoder.Token()
	if err != nil {
		syntaxErr := &xml.SyntaxError{}
		if err == io.EOF || xerrors.As(err, &syntaxErr) {
			finfo, _ := d.file.Stat()
			d.savedPos = finfo.Size()
			if err := d.finish(); err != nil {
				return xerrors.Errorf("failed to finish: %w", err)
			}
			return io.EOF
		}
		return xerrors.Errorf("failed to decode xml: %w", err)
	}
	start, ok := token.(xml.StartElement)
	if !ok {
		return nil
	}

	switch start.Name.Local {
	case "node":
		e, err := d.decodeNodeElement(&start)
		if err != nil {
			return xerrors.Errorf("failed to decode node: %w", err)
		}
		d.savedPos = offset
		d.lastElement = e
	case "way":
		e, err := d.decodeWayElement(&start)
		if err != nil {
			return xerrors.Errorf("failed to decode way: %w", err)
		}
		d.savedPos = offset
		d.lastElement = e
	case "relation":
		e, err := d.decodeRelationElement(&start)
		if err != nil {
			return xerrors.Errorf("failed to decode relation: %w", err)
		}
		d.savedPos = offset
		d.lastElement = e
	case "nd":
		i, _ := strconv.ParseInt(start.Attr[0].Value, 10, 64)
		w := d.lastElement.(*format.Way)
		w.NodeRefs = append(w.NodeRefs, &format.NodeRef{Ref: i})
	case "member":
		member := d.decodeMember(&start)
		r := d.lastElement.(*format.Relation)
		r.Members = append(r.Members, member)
	case "tag":
		d.lastElement.AddTag(start.Attr[0].Value, start.Attr[1].Value)
	}
	return nil
}

func (d *Decoder) Node(cb func(*format.Node) error) {
	d.nodeCallback = cb
}

func (d *Decoder) Way(cb func(*format.Way) error) {
	d.wayCallback = cb
}

func (d *Decoder) Relation(cb func(*format.Relation) error) {
	d.relCallback = cb
}

func (d *Decoder) Decode() error {
	if err := d.prepare(); err != nil {
		return xerrors.Errorf("failed to prepare decoding: %w", err)
	}
	if d.nodeCallback != nil {
		if err := d.storage.AllNodes(d.nodeCallback); err != nil {
			return xerrors.Errorf("failed to iterate node: %w", err)
		}
	}
	if d.wayCallback != nil {
		if err := d.storage.AllWays(d.wayCallback); err != nil {
			return xerrors.Errorf("failed to iterate way: %w", err)
		}
	}
	if d.relCallback != nil {
		if err := d.storage.AllRelations(d.relCallback); err != nil {
			return xerrors.Errorf("failed to iterate relation: %w", err)
		}
	}
	return nil
}
