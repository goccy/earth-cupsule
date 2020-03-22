package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/goccy/earth-cupsule/encoding/geojson"
	"github.com/goccy/earth-cupsule/encoding/osm"
	"github.com/goccy/earth-cupsule/format"
	"golang.org/x/xerrors"
)

var (
	ignoreNode     = flag.Bool("nonode", false, "no output from node")
	ignoreWay      = flag.Bool("noway", false, "no output from way")
	ignoreRelation = flag.Bool("norelation", false, "no output from relation")
)

type converter struct {
	dec *osm.Decoder
}

func (c *converter) close() {
	c.dec.Stop()
}

func (c *converter) convert(source, target string) (e error) {
	enc, err := geojson.NewEncoder(target)
	if err != nil {
		return xerrors.Errorf("failed to create geosjon.Encoder: %w", err)
	}
	defer func() {
		if err := enc.Close(); err != nil {
			e = xerrors.Errorf("failed to close geojson.Encoder: %w", err)
		}
	}()

	dec, err := osm.NewDecoder(source)
	if err != nil {
		return xerrors.Errorf("failed to create osm.Decoder: %w", err)
	}
	c.dec = dec
	defer func() {
		if err := dec.Close(); err != nil {
			e = xerrors.Errorf("failed to close osm.Decoder: %w", err)
		}
	}()
	if !*ignoreNode {
		dec.Node(func(node *format.Node) error {
			if err := enc.Encode(node.ToFeature()); err != nil {
				return xerrors.Errorf("failed to encode: %w", err)
			}
			return nil
		})
	}
	if !*ignoreWay {
		dec.Way(func(way *format.Way) error {
			if err := enc.Encode(way.ToFeature()); err != nil {
				return xerrors.Errorf("failed to encode: %w", err)
			}
			return nil
		})
	}
	if !*ignoreRelation {
		dec.Relation(func(rel *format.Relation) error {
			if err := enc.Encode(rel.ToFeature()); err != nil {
				return xerrors.Errorf("failed to encode: %w", err)
			}
			return nil
		})
	}
	if err := dec.Decode(); err != nil {
		return xerrors.Errorf("failed to decode: %w", err)
	}
	return nil
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		fmt.Println("osm2geojson must specify source.osm target.geojson")
		return
	}
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	var c converter
	go func() {
		for {
			s := <-sig
			switch s {
			case syscall.SIGINT:
				c.close()
			}
		}
	}()
	if err := c.convert(args[0], args[1]); err != nil {
		if xerrors.Is(err, osm.ErrForceStop) {
			return
		}
		log.Printf("%+v", err)
	}
}
