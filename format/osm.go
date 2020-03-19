package format

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type OSMType int

const (
	OSMNode     OSMType = 1
	OSMWay      OSMType = 2
	OSMRelation OSMType = 3
)

type OSM interface {
	Type() OSMType
	ToFeature() *Feature
	AddTag(string, string)
}

func NewNode() *Node {
	return &Node{
		Tags: Tags{},
	}
}

func NewWay() *Way {
	return &Way{
		Tags:     Tags{},
		NodeRefs: NodeRefs{},
	}
}

func NewRelation() *Relation {
	return &Relation{
		Tags:    Tags{},
		Members: Members{},
	}
}

type Tags map[string]string

// UninterestingTags are boring tags. If an element only has
// these tags it does not usually need to be displayed.
// For example, if a node with just these tags is part of a way, it
// probably does not need its own icon along the way.
var UninterestingTags = map[string]bool{
	"source":            true,
	"source_ref":        true,
	"source:ref":        true,
	"history":           true,
	"attribution":       true,
	"created_by":        true,
	"tiger:county":      true,
	"tiger:tlid":        true,
	"tiger:upload_uuid": true,
}

func (t Tags) Find(key string) string {
	return t[key]
}

func (t Tags) HasInterestingTags() bool {
	if len(t) == 0 {
		return false
	}

	for k := range t {
		if !UninterestingTags[k] {
			return true
		}
	}
	return false
}

func (t Tags) HasInterestingTagsWithTags(tags Tags) bool {
	if len(tags) == 0 {
		return false
	}
	for k, v := range t {
		if UninterestingTags[k] {
			continue
		}
		if tags.Find(k) == "true" ||
			tags.Find(k) == v {
			return true
		}
	}
	return false
}

type Node struct {
	OSM
	ID          int64     `json:"id"`
	Lat         float64   `json:"lat"`
	Lon         float64   `json:"lon"`
	User        string    `json:"user"`
	UserID      int64     `json:"uid"`
	Visible     bool      `json:"visible"`
	Version     int64     `json:"version"`
	ChangesetID int64     `json:"changeset"`
	Timestamp   time.Time `json:"timestamp"`
	Tags        Tags      `json:"tags"`

	buf *buffer
}

func (n *Node) Type() OSMType {
	return OSMNode
}

func (n *Node) ToFeature() *Feature {
	return NewFeature(
		fmt.Sprintf("node/%d", n.ID),
		Point{Lat: n.Lat, Lon: n.Lon},
		n.Tags,
	)
}

func (n *Node) AddTag(k, v string) {
	n.Tags[k] = v
}

func (n *Node) FeatureID() string {
	return fmt.Sprintf("node/%d", n.ID)
}

type buffer struct {
	b []byte
}

func newBuffer() *buffer {
	return &buffer{b: make([]byte, 0, 1024)}
}

func (b *buffer) reset() {
	b.b = b.b[:0]
}

func (b *buffer) appendByte(_b byte) {
	b.b = append(b.b, _b)
}

func (b *buffer) appendString(s string) {
	b.b = append(b.b, s...)
}

func (b *buffer) appendInt(i int64) {
	b.b = strconv.AppendInt(b.b, i, 10)
}

func (b *buffer) appendFloat(f float64) {
	b.b = strconv.AppendFloat(b.b, f, 'E', -1, 64)
}

func (b *buffer) appendBool(_b bool) {
	b.b = strconv.AppendBool(b.b, _b)
}

var bufferpool = sync.Pool{
	New: func() interface{} {
		return newBuffer()
	},
}

func (n *Node) Release() {
	n.buf.reset()
	bufferpool.Put(n.buf)
}

func (n *Node) MarshalJSON() ([]byte, error) {
	buf := bufferpool.Get().(*buffer)
	buf.appendByte('{')
	buf.appendString(`"id":`)
	buf.appendInt(n.ID)
	buf.appendByte(',')
	buf.appendString(`"lat":`)
	buf.appendFloat(n.Lat)
	buf.appendByte(',')
	buf.appendString(`"lon":`)
	buf.appendFloat(n.Lon)
	buf.appendByte(',')
	buf.appendString(`"user":`)
	buf.appendString(fmt.Sprintf(`"%s"`, n.User))
	buf.appendByte(',')
	buf.appendString(`"userId":`)
	buf.appendInt(n.UserID)
	buf.appendByte(',')
	buf.appendString(`"visible":`)
	buf.appendBool(n.Visible)
	buf.appendByte(',')
	buf.appendString(`"version":`)
	buf.appendInt(int64(n.Version))
	buf.appendByte(',')
	buf.appendString(`"changesetId":`)
	buf.appendInt(n.ChangesetID)
	buf.appendByte(',')
	buf.appendString(`"timestamp":`)
	buf.appendString(fmt.Sprintf(`"%s"`, n.Timestamp.Format("2006-01-02T15:04:05.9Z")))
	buf.appendByte(',')
	buf.appendString(`"tags":{`)
	i := 0
	for k, v := range n.Tags {
		if i > 0 {
			buf.appendByte(',')
		}
		buf.appendString(fmt.Sprintf(`"%s":%s`, k, strconv.Quote(v)))
		i++
	}
	buf.appendByte('}')
	buf.appendByte('}')
	n.buf = buf
	return buf.b, nil
}

func (n *Node) GetID() int64 {
	return n.ID
}

type NodeRef struct {
	Ref int64
}

func (n *NodeRef) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.Ref)
}

func (n *NodeRef) UnmarshalJSON(bytes []byte) error {
	var v int64
	if err := json.Unmarshal(bytes, &v); err != nil {
		return errors.WithStack(err)
	}
	n.Ref = v
	return nil
}

type NodeRefs []*NodeRef

type Way struct {
	OSM
	ID          int64
	User        string
	UserID      int64
	Visible     bool
	Version     int64
	ChangesetID int64
	Timestamp   time.Time
	NodeRefs    NodeRefs
	Tags        Tags
	Skippable   bool `xml:"-"`

	Nodes []*Node `xml:"-"`
}

func (w *Way) toLineString() LineString {
	points := []Point{}
	for _, node := range w.Nodes {
		points = append(points, Point{
			Lat: node.Lat,
			Lon: node.Lon,
		})
	}
	return LineString(points)
}

func (w *Way) ToFeature() *Feature {
	return NewFeature(fmt.Sprintf("way/%d", w.ID), w.toLineString(), w.Tags)
}

func (w *Way) Type() OSMType {
	return OSMWay
}

func (w *Way) AddTag(k, v string) {
	w.Tags[k] = v
}

func (w *Way) FeatureID() string {
	return fmt.Sprintf("way/%d", w.ID)
}

func (w *Way) NodeIDs() []int64 {
	ids := make([]int64, len(w.NodeRefs))
	for idx, nodeRef := range w.NodeRefs {
		ids[idx] = nodeRef.Ref
	}
	return ids
}

func (w *Way) GetID() int64 {
	return w.ID
}

type Relation struct {
	OSM
	ID          int64
	User        string
	UserID      int64
	Visible     bool
	Version     int64
	ChangesetID int64
	Timestamp   time.Time
	Tags        Tags
	Members     Members
}

func (r *Relation) ToFeature() *Feature {
	switch r.Tags.Find("type") {
	case "route":
		return r.routeToFeature()
	case "multipolygon", "boundary":
		return r.polygonToFeature()
	}
	return nil
}

func (r *Relation) routeToFeature() *Feature {
	lines := []LineString{}
	for _, m := range r.Members {
		if !m.IsWay() {
			continue
		}
		line := m.Way.toLineString()
		if len(line) == 0 {
			continue
		}
		lines = append(lines, line)
	}
	if len(lines) == 0 {
		return nil
	}
	if len(lines) == 1 {
		return NewFeature(fmt.Sprintf("relation/%d", r.ID), lines[0], r.Tags)
	}
	ml := MultiLineString(lines)
	return NewFeature(fmt.Sprintf("relation/%d", r.ID), ml, r.Tags)
}

func (r *Relation) polygonToFeature() *Feature {
	inners := r.Members.Inners()
	outers := r.Members.Outers()
	innerLines := []LineString{}
	outerLines := []LineString{}
	for _, inner := range inners {
		innerLines = append(innerLines, inner.toLineString())
	}
	for _, outer := range outers {
		line := outer.toLineString()
		if len(line) == 0 {
			continue
		}
		outerLines = append(outerLines, line)
	}
	if len(outerLines) == 0 {
		return nil
	}
	if len(outerLines) == 1 {
		polygon := Polygon(outerLines)
		polygon = append(polygon, innerLines...)
		if !r.Tags.HasInterestingTagsWithTags(Tags(map[string]string{"type": "true"})) {
			// skippable
		}
		return NewFeature(fmt.Sprintf("relation/%d", r.ID), polygon, r.Tags)
	}
	return nil
}

func (r *Relation) Type() OSMType {
	return OSMRelation
}

func (r *Relation) AddTag(k, v string) {
	r.Tags[k] = v
}

func (r *Relation) FeatureID() string {
	return fmt.Sprintf("relation/%d", r.ID)
}

func (r *Relation) UnmarshalRelation(bytes []byte) error {
	return json.Unmarshal(bytes, r)
}

func (r *Relation) GetID() int64 {
	return r.ID
}

func (r *Relation) NodeTypeMemberIDs() []int64 {
	ids := []int64{}
	for _, member := range r.Members {
		if member.Type != TypeNode {
			continue
		}
		ids = append(ids, member.Ref)
	}
	return ids
}

type Orientation int

const (
	OrientationCCW Orientation = 1  // CounerClockWise
	OrientationCW              = -1 // ClockWise
)

type Type string

const (
	TypeNode      Type = "node"
	TypeWay            = "way"
	TypeRelation       = "relation"
	TypeChangeset      = "changeset"
	TypeNote           = "note"
	TypeUser           = "user"
)

type Member struct {
	Type        Type
	Ref         int64
	Role        string
	Version     int64
	ChangesetID int64
	Lat         float64
	Lon         float64
	Orientation Orientation

	Way *Way `xml:"-"`
}

func (m *Member) IsWay() bool {
	if m.Way == nil {
		return false
	}
	if m.Way.ID == 0 {
		return false
	}
	return m.Type == TypeWay
}

type Members []*Member

func (m Members) Inners() []*Way {
	inners := []*Way{}
	for _, mm := range m {
		if !mm.IsWay() {
			continue
		}
		if mm.Role != "inner" {
			continue
		}
		inners = append(inners, mm.Way)
	}
	return inners
}

func (m Members) Outers() []*Way {
	outers := []*Way{}
	for _, mm := range m {
		if !mm.IsWay() {
			continue
		}
		if mm.Role != "outer" {
			continue
		}
		outers = append(outers, mm.Way)
	}
	return outers
}
