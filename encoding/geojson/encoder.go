package geojson

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/goccy/earth-cupsule/format"
	"golang.org/x/xerrors"
)

type Encoder struct {
	file           *os.File
	jsonEncoder    *json.Encoder
	isFirstFeature bool
	features       []*format.Feature
	mu             sync.Mutex
}

const (
	flushThreshold = 10000
)

var (
	header = []byte(`{"type": "FeatureCollection","features": [`)
	footer = []byte(`]}`)
)

func NewEncoder(path string) (*Encoder, error) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, xerrors.Errorf("failed to open file %s: %w", path, err)
	}
	return &Encoder{
		file:           f,
		jsonEncoder:    json.NewEncoder(f),
		isFirstFeature: true,
		features:       make([]*format.Feature, 0, flushThreshold),
	}, nil
}

func (e *Encoder) Close() error {
	if err := e.flush(); err != nil {
		return xerrors.Errorf("failed to flush: %w", err)
	}
	if _, err := e.file.Write(footer); err != nil {
		return xerrors.Errorf("failed to write footer: %w", err)
	}
	if err := e.file.Close(); err != nil {
		return xerrors.Errorf("failed to close: %w", err)
	}
	return nil
}

func (e *Encoder) flush() error {
	for _, feature := range e.features {
		if e.isFirstFeature {
			if _, err := e.file.Write(header); err != nil {
				return xerrors.Errorf("failed to write header: %w", err)
			}
			e.isFirstFeature = false
		} else {
			if _, err := e.file.Write([]byte(`,`)); err != nil {
				return xerrors.Errorf("failed to write delim: %w", err)
			}
		}
		if err := e.jsonEncoder.Encode(feature); err != nil {
			return xerrors.Errorf("failed to write feature: %w", err)
		}
	}
	e.features = e.features[:0]
	return nil
}

func (e *Encoder) addFeature(feature *format.Feature) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.features = append(e.features, feature)
	if len(e.features) > flushThreshold {
		if err := e.flush(); err != nil {
			return xerrors.Errorf("failed to flush: %w", err)
		}
	}
	return nil
}

func (e *Encoder) Encode(feature *format.Feature) error {
	if err := e.addFeature(feature); err != nil {
		return xerrors.Errorf("failed to addFeature: %w", err)
	}
	return nil
}
