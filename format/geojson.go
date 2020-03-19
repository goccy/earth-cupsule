package format

import "encoding/json"

type Feature struct {
	ID         string
	Geometry   Geometry
	Properties map[string]string
}

func NewFeature(id string, geometry Geometry, props map[string]string) *Feature {
	return &Feature{
		ID:         id,
		Geometry:   geometry,
		Properties: props,
	}
}

func (f *Feature) MarshalJSON() ([]byte, error) {
	geometry := struct {
		Type        string      `json:"type"`
		Coordinates interface{} `json:"coordinates"`
	}{
		Type:        f.Geometry.Type(),
		Coordinates: f.Geometry.Coordinates(),
	}
	return json.Marshal(struct {
		Type       string            `json:"type"`
		ID         string            `json:"id"`
		Geometry   interface{}       `json:"geometry"`
		Properties map[string]string `json:"properties"`
	}{
		Type:       "Feature",
		ID:         f.ID,
		Geometry:   geometry,
		Properties: f.Properties,
	})
}

type Geometry interface {
	Type() string
	Coordinates() interface{}
}

type Point struct {
	Lat float64
	Lon float64
}

func (p Point) Type() string {
	return "Point"
}

func (p Point) Coordinates() interface{} {
	return []float64{p.Lon, p.Lat}
}

type MultiPoint []Point

func (m MultiPoint) Type() string {
	return "MultiPoint"
}

func (m MultiPoint) Coordinates() interface{} {
	coords := []interface{}{}
	for _, p := range m {
		coords = append(coords, p.Coordinates())
	}
	return coords
}

type LineString []Point

func (l LineString) Type() string {
	return "LineString"
}

func (l LineString) Coordinates() interface{} {
	coords := []interface{}{}
	for _, p := range l {
		coords = append(coords, p.Coordinates())
	}
	return coords
}

type MultiLineString []LineString

func (m MultiLineString) Type() string {
	return "MultiLineString"
}

func (m MultiLineString) Coordinates() interface{} {
	coords := []interface{}{}
	for _, l := range m {
		coords = append(coords, l.Coordinates())
	}
	return coords
}

type Polygon []LineString

func (p Polygon) Type() string {
	return "Polygon"
}

func (p Polygon) Coordinates() interface{} {
	coords := []interface{}{}
	for _, l := range p {
		coords = append(coords, l.Coordinates())
	}
	return coords
}

type MultiPolygon []Polygon

func (m MultiPolygon) Type() string {
	return "MultiPolygon"
}

func (m MultiPolygon) Coordinates() interface{} {
	coords := []interface{}{}
	for _, p := range m {
		coords = append(coords, p.Coordinates())
	}
	return coords
}
