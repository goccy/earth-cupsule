package osm

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/goccy/earth-cupsule/format"
	"github.com/mitchellh/ioprogress"
	"golang.org/x/xerrors"
)

var (
	ErrForceStop = xerrors.New("force stop")
)

type Decoder struct {
	file         *os.File
	storage      *Storage
	dec          StreamDecoder
	savedPos     int64
	isNeededStop bool
	onceClose    sync.Once
	nodeCallback func(*format.Node) error
	wayCallback  func(*format.Way) error
	relCallback  func(*format.Relation) error
}

type StreamDecoder interface {
	IsDecoded() bool
	Decode() (int64, error)
	SetNodeCallback(func(*format.Node) error)
	Close()
}

func NewDecoder(path string) (*Decoder, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, xerrors.Errorf("failed to open file: %w", err)
	}
	s, err := NewStorage(path)
	if err != nil {
		return nil, xerrors.Errorf("failed to open storage for osm: %w", err)
	}
	pos := s.Pos()
	if _, err := f.Seek(pos, os.SEEK_SET); err != nil {
		return nil, xerrors.Errorf("failed to seek: %w", err)
	}
	finfo, err := f.Stat()
	if err != nil {
		return nil, xerrors.Errorf("failed to get file stat: %w", err)
	}
	fsize := finfo.Size()
	var progress *ioprogress.Reader
	if pos < finfo.Size() {
		progress = &ioprogress.Reader{
			Reader: f,
			Size:   fsize,
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
	var dec StreamDecoder
	if filepath.Ext(path) == ".pbf" {
		dec = NewPBFDecoder(progress, s, fsize)
	} else {
		dec = NewXMLDecoder(progress, s, fsize)
	}
	return &Decoder{
		file:    f,
		storage: s,
		dec:     dec,
	}, nil
}

func (d *Decoder) close() error {
	d.dec.Close()
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

func (d *Decoder) callbackDecodedNodes() error {
	if d.nodeCallback == nil {
		return nil
	}
	if err := d.storage.AllNodes(func(node *format.Node) error {
		if d.isNeededStop {
			return ErrForceStop
		}
		if err := d.nodeCallback(node); err != nil {
			return xerrors.Errorf("failed to node callback: %w", err)
		}
		return nil
	}); err != nil {
		return xerrors.Errorf("failed to iterate node: %w", err)
	}
	return nil
}

func (d *Decoder) prepare() error {
	if err := d.callbackDecodedNodes(); err != nil {
		return xerrors.Errorf("failed to callback for decoded node: %w", err)
	}
	if d.dec.IsDecoded() {
		return nil
	}
	for {
		if d.isNeededStop {
			return ErrForceStop
		}
		pos, err := d.dec.Decode()
		if err != nil {
			if xerrors.Is(err, io.EOF) {
				d.savedPos = pos
				if err := d.finish(); err != nil {
					return xerrors.Errorf("failed to finish: %w", err)
				}
				break
			}
			return xerrors.Errorf("failed to decode: %w", err)
		}
		if pos > 0 {
			d.savedPos = pos
		}
	}
	return nil
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

func (d *Decoder) Node(cb func(*format.Node) error) {
	d.dec.SetNodeCallback(cb)
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
	if d.wayCallback != nil {
		if err := d.storage.AllWays(func(way *format.Way) error {
			if d.isNeededStop {
				return ErrForceStop
			}
			if err := d.wayCallback(way); err != nil {
				return xerrors.Errorf("failed to way callback: %w", err)
			}
			return nil
		}); err != nil {
			return xerrors.Errorf("failed to iterate way: %w", err)
		}
	}
	if d.relCallback != nil {
		if err := d.storage.AllRelations(func(rel *format.Relation) error {
			if d.isNeededStop {
				return ErrForceStop
			}
			if err := d.relCallback(rel); err != nil {
				return xerrors.Errorf("failed to rel callback: %w", err)
			}
			return nil
		}); err != nil {
			return xerrors.Errorf("failed to iterate relation: %w", err)
		}
	}
	return nil
}
