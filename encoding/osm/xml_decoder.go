package osm

import (
	"encoding/xml"
	"io"
	"strconv"
	"time"

	"github.com/goccy/earth-cupsule/format"
	"golang.org/x/xerrors"
)

type XMLDecoder struct {
	fsize       int64
	storage     *Storage
	xmlDecoder  *xml.Decoder
	lastElement format.OSM
	offset      int64
}

func NewXMLDecoder(reader io.Reader, storage *Storage, fsize int64) *XMLDecoder {
	return &XMLDecoder{
		fsize:      fsize,
		storage:    storage,
		xmlDecoder: xml.NewDecoder(reader),
		offset:     storage.Pos(),
	}
}

func (d *XMLDecoder) IsDecoded() bool {
	return d.offset >= d.fsize
}

func (d *XMLDecoder) setNextElement(e format.OSM) error {
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

func (d *XMLDecoder) decodeNodeElement(e *xml.StartElement) (*format.Node, error) {
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

func (d *XMLDecoder) decodeWayElement(e *xml.StartElement) (*format.Way, error) {
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

func (d *XMLDecoder) decodeRelationElement(e *xml.StartElement) (*format.Relation, error) {
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

func (d *XMLDecoder) decodeMember(e *xml.StartElement) *format.Member {
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

func (d *XMLDecoder) Decode() (int64, error) {
	offset := d.offset + d.xmlDecoder.InputOffset()
	token, err := d.xmlDecoder.Token()
	if err != nil {
		syntaxErr := &xml.SyntaxError{}
		if err == io.EOF || xerrors.As(err, &syntaxErr) {
			return d.fsize, io.EOF
		}
		return 0, xerrors.Errorf("failed to decode xml: %w", err)
	}
	start, ok := token.(xml.StartElement)
	if !ok {
		return 0, nil
	}

	switch start.Name.Local {
	case "node":
		e, err := d.decodeNodeElement(&start)
		if err != nil {
			return 0, xerrors.Errorf("failed to decode node: %w", err)
		}
		d.lastElement = e
		return offset, nil
	case "way":
		e, err := d.decodeWayElement(&start)
		if err != nil {
			return 0, xerrors.Errorf("failed to decode way: %w", err)
		}
		d.lastElement = e
		return offset, nil
	case "relation":
		e, err := d.decodeRelationElement(&start)
		if err != nil {
			return 0, xerrors.Errorf("failed to decode relation: %w", err)
		}
		d.lastElement = e
		return offset, nil
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
	return 0, nil
}
