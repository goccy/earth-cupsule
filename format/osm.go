package format

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/vmihailenco/msgpack/v4"
	"golang.org/x/xerrors"
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

func (n *Node) EncodeMsgpack(enc *msgpack.Encoder) error {
	if err := enc.EncodeInt64(n.ID); err != nil {
		return xerrors.Errorf("failed to encode id: %w", err)
	}
	if err := enc.EncodeFloat64(n.Lat); err != nil {
		return xerrors.Errorf("failed to encode lat: %w", err)
	}
	if err := enc.EncodeFloat64(n.Lon); err != nil {
		return xerrors.Errorf("failed to encode lon: %w", err)
	}
	if err := enc.EncodeString(n.User); err != nil {
		return xerrors.Errorf("failed to encode user: %w", err)
	}
	if err := enc.EncodeInt64(n.UserID); err != nil {
		return xerrors.Errorf("failed to encode user_id: %w", err)
	}
	if err := enc.EncodeBool(n.Visible); err != nil {
		return xerrors.Errorf("failed to encode visible: %w", err)
	}
	if err := enc.EncodeInt64(n.Version); err != nil {
		return xerrors.Errorf("failed to encode version: %w", err)
	}
	if err := enc.EncodeInt64(n.ChangesetID); err != nil {
		return xerrors.Errorf("failed to encode changeset: %w", err)
	}
	if err := enc.EncodeTime(n.Timestamp); err != nil {
		return xerrors.Errorf("failed to encode timestamp: %w", err)
	}
	if err := enc.EncodeMapLen(len(n.Tags)); err != nil {
		return xerrors.Errorf("failed to encode tags length: %w", err)
	}
	for k, v := range n.Tags {
		if err := enc.EncodeString(k); err != nil {
			return xerrors.Errorf("failed to encode key: %w", err)
		}
		if err := enc.EncodeString(v); err != nil {
			return xerrors.Errorf("failed to encode value: %w", err)
		}
	}
	return nil
}

func (n *Node) DecodeMsgpack(dec *msgpack.Decoder) error {
	id, err := dec.DecodeInt64()
	if err != nil {
		return xerrors.Errorf("failed to decode id: %w", err)
	}
	lat, err := dec.DecodeFloat64()
	if err != nil {
		return xerrors.Errorf("failed to decode lat: %w", err)
	}
	lon, err := dec.DecodeFloat64()
	if err != nil {
		return xerrors.Errorf("failed to decode lon: %w", err)
	}
	user, err := dec.DecodeString()
	if err != nil {
		return xerrors.Errorf("failed to decode user: %w", err)
	}
	userID, err := dec.DecodeInt64()
	if err != nil {
		return xerrors.Errorf("failed to decode user_id: %w", err)
	}
	visible, err := dec.DecodeBool()
	if err != nil {
		return xerrors.Errorf("failed to decode visible: %w", err)
	}
	version, err := dec.DecodeInt64()
	if err != nil {
		return xerrors.Errorf("failed to decode version: %w", err)
	}
	changeset, err := dec.DecodeInt64()
	if err != nil {
		return xerrors.Errorf("failed to decode changeset: %w", err)
	}
	timestamp, err := dec.DecodeTime()
	if err != nil {
		return xerrors.Errorf("failed to decode timestamp: %w", err)
	}
	tagLen, err := dec.DecodeMapLen()
	if err != nil {
		return xerrors.Errorf("failed to decode tags length: %w", err)
	}
	tags := make(Tags, tagLen)
	for i := 0; i < tagLen; i++ {
		key, err := dec.DecodeString()
		if err != nil {
			return xerrors.Errorf("failed to decode tag key: %w", err)
		}
		value, err := dec.DecodeString()
		if err != nil {
			return xerrors.Errorf("failed to decode tag value: %w", err)
		}
		tags[key] = value
	}
	n.ID = id
	n.Lat = lat
	n.Lon = lon
	n.User = user
	n.UserID = userID
	n.Visible = visible
	n.Version = version
	n.ChangesetID = changeset
	n.Timestamp = timestamp
	n.Tags = tags
	return nil
}

func (n *Node) GetID() int64 {
	return n.ID
}

type NodeRef struct {
	Ref int64
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

func (w *Way) EncodeMsgpack(enc *msgpack.Encoder) error {
	if err := enc.EncodeInt64(w.ID); err != nil {
		return xerrors.Errorf("failed to encode id: %w", err)
	}
	if err := enc.EncodeString(w.User); err != nil {
		return xerrors.Errorf("failed to encode user: %w", err)
	}
	if err := enc.EncodeInt64(w.UserID); err != nil {
		return xerrors.Errorf("failed to encode user_id: %w", err)
	}
	if err := enc.EncodeBool(w.Visible); err != nil {
		return xerrors.Errorf("failed to encode visible: %w", err)
	}
	if err := enc.EncodeInt64(w.Version); err != nil {
		return xerrors.Errorf("failed to encode version: %w", err)
	}
	if err := enc.EncodeInt64(w.ChangesetID); err != nil {
		return xerrors.Errorf("failed to encode changeset: %w", err)
	}
	if err := enc.EncodeTime(w.Timestamp); err != nil {
		return xerrors.Errorf("failed to encode timestamp: %w", err)
	}
	if err := enc.EncodeArrayLen(len(w.NodeRefs)); err != nil {
		return xerrors.Errorf("failed to encode noderefs length: %w", err)
	}
	for _, ref := range w.NodeRefs {
		if err := enc.EncodeInt64(ref.Ref); err != nil {
			return xerrors.Errorf("failed to encode ref: %w", err)
		}
	}
	if err := enc.EncodeMapLen(len(w.Tags)); err != nil {
		return xerrors.Errorf("failed to encode tags length: %w", err)
	}
	for k, v := range w.Tags {
		if err := enc.EncodeString(k); err != nil {
			return xerrors.Errorf("failed to encode key: %w", err)
		}
		if err := enc.EncodeString(v); err != nil {
			return xerrors.Errorf("failed to encode value: %w", err)
		}
	}
	return nil
}

func (w *Way) DecodeMsgpack(dec *msgpack.Decoder) error {
	id, err := dec.DecodeInt64()
	if err != nil {
		return xerrors.Errorf("failed to decode id: %w", err)
	}
	user, err := dec.DecodeString()
	if err != nil {
		return xerrors.Errorf("failed to decode user: %w", err)
	}
	userID, err := dec.DecodeInt64()
	if err != nil {
		return xerrors.Errorf("failed to decode user_id: %w", err)
	}
	visible, err := dec.DecodeBool()
	if err != nil {
		return xerrors.Errorf("failed to decode visible: %w", err)
	}
	version, err := dec.DecodeInt64()
	if err != nil {
		return xerrors.Errorf("failed to decode version: %w", err)
	}
	changeset, err := dec.DecodeInt64()
	if err != nil {
		return xerrors.Errorf("failed to decode changeset: %w", err)
	}
	timestamp, err := dec.DecodeTime()
	if err != nil {
		return xerrors.Errorf("failed to decode timestamp: %w", err)
	}
	nodeRefLen, err := dec.DecodeArrayLen()
	if err != nil {
		return xerrors.Errorf("failed to decode noderefs length: %w", err)
	}
	nodeRefs := make(NodeRefs, nodeRefLen)
	for i := 0; i < nodeRefLen; i++ {
		ref, err := dec.DecodeInt64()
		if err != nil {
			return xerrors.Errorf("failed to decode ref: %w", err)
		}
		nodeRefs[i] = &NodeRef{Ref: ref}
	}
	tagLen, err := dec.DecodeMapLen()
	if err != nil {
		return xerrors.Errorf("failed to decode tags length: %w", err)
	}
	tags := make(Tags, tagLen)
	for i := 0; i < tagLen; i++ {
		key, err := dec.DecodeString()
		if err != nil {
			return xerrors.Errorf("failed to decode tag key: %w", err)
		}
		value, err := dec.DecodeString()
		if err != nil {
			return xerrors.Errorf("failed to decode tag value: %w", err)
		}
		tags[key] = value
	}
	w.ID = id
	w.User = user
	w.UserID = userID
	w.Visible = visible
	w.Version = version
	w.ChangesetID = changeset
	w.Timestamp = timestamp
	w.NodeRefs = nodeRefs
	w.Tags = tags
	return nil
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

func (r *Relation) EncodeMsgpack(enc *msgpack.Encoder) error {
	if err := enc.EncodeInt64(r.ID); err != nil {
		return xerrors.Errorf("failed to encode id: %w", err)
	}
	if err := enc.EncodeString(r.User); err != nil {
		return xerrors.Errorf("failed to encode user: %w", err)
	}
	if err := enc.EncodeInt64(r.UserID); err != nil {
		return xerrors.Errorf("failed to encode user_id: %w", err)
	}
	if err := enc.EncodeBool(r.Visible); err != nil {
		return xerrors.Errorf("failed to encode visible: %w", err)
	}
	if err := enc.EncodeInt64(r.Version); err != nil {
		return xerrors.Errorf("failed to encode version: %w", err)
	}
	if err := enc.EncodeInt64(r.ChangesetID); err != nil {
		return xerrors.Errorf("failed to encode changeset: %w", err)
	}
	if err := enc.EncodeTime(r.Timestamp); err != nil {
		return xerrors.Errorf("failed to encode timestamp: %w", err)
	}
	if err := enc.EncodeArrayLen(len(r.Members)); err != nil {
		return xerrors.Errorf("failed to encode members length: %w", err)
	}
	for _, member := range r.Members {
		if err := member.EncodeMsgpack(enc); err != nil {
			return xerrors.Errorf("failed to encode member: %w", err)
		}
	}
	if err := enc.EncodeMapLen(len(r.Tags)); err != nil {
		return xerrors.Errorf("failed to encode tags length: %w", err)
	}
	for k, v := range r.Tags {
		if err := enc.EncodeString(k); err != nil {
			return xerrors.Errorf("failed to encode key: %w", err)
		}
		if err := enc.EncodeString(v); err != nil {
			return xerrors.Errorf("failed to encode value: %w", err)
		}
	}
	return nil
}

func (r *Relation) DecodeMsgpack(dec *msgpack.Decoder) error {
	id, err := dec.DecodeInt64()
	if err != nil {
		return xerrors.Errorf("failed to decode id: %w", err)
	}
	user, err := dec.DecodeString()
	if err != nil {
		return xerrors.Errorf("failed to decode user: %w", err)
	}
	userID, err := dec.DecodeInt64()
	if err != nil {
		return xerrors.Errorf("failed to decode user_id: %w", err)
	}
	visible, err := dec.DecodeBool()
	if err != nil {
		return xerrors.Errorf("failed to decode visible: %w", err)
	}
	version, err := dec.DecodeInt64()
	if err != nil {
		return xerrors.Errorf("failed to decode version: %w", err)
	}
	changeset, err := dec.DecodeInt64()
	if err != nil {
		return xerrors.Errorf("failed to decode changeset: %w", err)
	}
	timestamp, err := dec.DecodeTime()
	if err != nil {
		return xerrors.Errorf("failed to decode timestamp: %w", err)
	}
	membersLen, err := dec.DecodeArrayLen()
	if err != nil {
		return xerrors.Errorf("failed to decode members length: %w", err)
	}
	members := make(Members, membersLen)
	for i := 0; i < membersLen; i++ {
		member := new(Member)
		if err := member.DecodeMsgpack(dec); err != nil {
			return xerrors.Errorf("failed to decode member: %w", err)
		}
		members[i] = member
	}
	tagLen, err := dec.DecodeMapLen()
	if err != nil {
		return xerrors.Errorf("failed to decode tags length: %w", err)
	}
	tags := make(Tags, tagLen)
	for i := 0; i < tagLen; i++ {
		key, err := dec.DecodeString()
		if err != nil {
			return xerrors.Errorf("failed to decode tag key: %w", err)
		}
		value, err := dec.DecodeString()
		if err != nil {
			return xerrors.Errorf("failed to decode tag value: %w", err)
		}
		tags[key] = value
	}
	r.ID = id
	r.User = user
	r.UserID = userID
	r.Visible = visible
	r.Version = version
	r.ChangesetID = changeset
	r.Timestamp = timestamp
	r.Members = members
	r.Tags = tags
	return nil
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

func (m *Member) EncodeMsgpack(enc *msgpack.Encoder) error {
	if err := enc.EncodeString(string(m.Type)); err != nil {
		return xerrors.Errorf("failed to encode type: %w", err)
	}
	if err := enc.EncodeInt64(m.Ref); err != nil {
		return xerrors.Errorf("failed to encode ref: %w", err)
	}
	if err := enc.EncodeString(m.Role); err != nil {
		return xerrors.Errorf("failed to encode role: %w", err)
	}
	if err := enc.EncodeInt64(m.Version); err != nil {
		return xerrors.Errorf("failed to encode version: %w", err)
	}
	if err := enc.EncodeInt64(m.ChangesetID); err != nil {
		return xerrors.Errorf("failed to encode changeset: %w", err)
	}
	if err := enc.EncodeFloat64(m.Lat); err != nil {
		return xerrors.Errorf("failed to encode lat: %w", err)
	}
	if err := enc.EncodeFloat64(m.Lon); err != nil {
		return xerrors.Errorf("failed to encode lon: %w", err)
	}
	if err := enc.EncodeInt8(int8(m.Orientation)); err != nil {
		return xerrors.Errorf("failed to encode orientation: %w", err)
	}
	return nil
}

func (m *Member) DecodeMsgpack(dec *msgpack.Decoder) error {
	typ, err := dec.DecodeString()
	if err != nil {
		return xerrors.Errorf("failed to decode type: %w", err)
	}
	ref, err := dec.DecodeInt64()
	if err != nil {
		return xerrors.Errorf("failed to decode id: %w", err)
	}
	role, err := dec.DecodeString()
	if err != nil {
		return xerrors.Errorf("failed to decode role: %w", err)
	}
	version, err := dec.DecodeInt64()
	if err != nil {
		return xerrors.Errorf("failed to decode version: %w", err)
	}
	changeset, err := dec.DecodeInt64()
	if err != nil {
		return xerrors.Errorf("failed to decode changeset: %w", err)
	}
	lat, err := dec.DecodeFloat64()
	if err != nil {
		return xerrors.Errorf("failed to decode lat: %w", err)
	}
	lon, err := dec.DecodeFloat64()
	if err != nil {
		return xerrors.Errorf("failed to decode lon: %w", err)
	}
	orientation, err := dec.DecodeInt8()
	if err != nil {
		return xerrors.Errorf("failed to decode orientation: %w", err)
	}
	m.Type = Type(typ)
	m.Ref = ref
	m.Role = role
	m.Version = version
	m.ChangesetID = changeset
	m.Lat = lat
	m.Lon = lon
	m.Orientation = Orientation(orientation)
	return nil
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
