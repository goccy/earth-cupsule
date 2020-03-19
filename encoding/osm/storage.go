package osm

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/dgraph-io/badger/v2"
	"github.com/dgraph-io/badger/v2/pb"
	"golang.org/x/xerrors"
)

type item struct {
	key   []byte
	value []byte
}

type Storage struct {
	db     *badger.DB
	items  []*item
	mu     sync.Mutex
	itemMu sync.Mutex
}

const (
	maxItemNum = 10000
)

type Logger struct{}

func (*Logger) Errorf(string, ...interface{})   {}
func (*Logger) Warningf(string, ...interface{}) {}
func (*Logger) Infof(string, ...interface{})    {}
func (*Logger) Debugf(string, ...interface{})   {}

func NewStorage(path string) (*Storage, error) {
	var logger Logger
	db, err := badger.Open(badger.DefaultOptions(filepath.Join("osmdb")).WithLogger(&logger))
	if err != nil {
		return nil, xerrors.Errorf("failed to open badger db for osm: %w", err)
	}
	return &Storage{
		db:    db,
		items: []*item{},
	}, nil
}

func (s *Storage) Close() error {
	if err := s.db.Close(); err != nil {
		return xerrors.Errorf("failed to close db for osm: %w", err)
	}
	return nil
}

func (s *Storage) toID(header string, id int64) []byte {
	return []byte(fmt.Sprintf("%s/%d", header, id))
}

func (s *Storage) nodeID(id int64) []byte {
	return s.toID("node", id)
}

func (s *Storage) nodeIDInWay(id int64) []byte {
	return s.toID("wnode", id)
}

func (s *Storage) wayIDInMember(id int64) []byte {
	return s.toID("mnode", id)
}

func (s *Storage) wayID(id int64) []byte {
	return s.toID("way", id)
}

func (s *Storage) relationID(id int64) []byte {
	return s.toID("rel", id)
}

func (s *Storage) getItemByKey(key []byte) ([]byte, error) {
	var b []byte
	if err := s.db.View(func(tx *badger.Txn) error {
		item, err := tx.Get(key)
		if err != nil {
			return xerrors.Errorf("failed to get item: %w", err)
		}
		v, err := item.ValueCopy(nil)
		if err != nil {
			return xerrors.Errorf("failed to get value: %w", err)
		}
		b = v
		return nil
	}); err != nil {
		return nil, xerrors.Errorf("failed to get item by key: %w", err)
	}
	return b, nil
}

func (s *Storage) Finish() error {
	if err := s.commitBufferedItems(); err != nil {
		return xerrors.Errorf("failed to commit buffered items: %w", err)
	}
	if err := s.db.Sync(); err != nil {
		return xerrors.Errorf("failed to sync: %w", err)
	}
	return nil
}

func (s *Storage) commitBufferedItems() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	tx := s.db.NewTransaction(true)
	if err := s.iterItems(func(item *item) error {
		if err := tx.SetEntry(badger.NewEntry(item.key, item.value)); err != nil {
			return xerrors.Errorf("failed to set k/v: %w", err)
		}
		return nil
	}); err != nil {
		return xerrors.Errorf("failed to items: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return xerrors.Errorf("failed to commit: %w", err)
	}
	s.clearItem()
	return nil
}

func (s *Storage) iterItems(cb func(item *item) error) error {
	s.itemMu.Lock()
	defer s.itemMu.Unlock()
	for _, item := range s.items {
		if err := cb(item); err != nil {
			return xerrors.Errorf("failed to callack item: %w", err)
		}
	}
	return nil
}

func (s *Storage) itemLen() int {
	s.itemMu.Lock()
	defer s.itemMu.Unlock()
	return len(s.items)
}

func (s *Storage) appendItem(key []byte, value []byte) {
	s.itemMu.Lock()
	defer s.itemMu.Unlock()
	v := make([]byte, len(value))
	copy(v, value)
	s.items = append(s.items, &item{key: key, value: v})
}

func (s *Storage) clearItem() {
	s.itemMu.Lock()
	defer s.itemMu.Unlock()
	s.items = []*item{}
}

func (s *Storage) setItem(key []byte, value []byte) error {
	if s.itemLen() > maxItemNum {
		if err := s.commitBufferedItems(); err != nil {
			return xerrors.Errorf("failed to commit buffered items: %w", err)
		}
	}
	s.appendItem(key, value)
	return nil
}

func (s *Storage) Pos() int64 {
	pos, err := s.getItemByKey([]byte("pos"))
	if err != nil {
		return 0
	}
	i, err := strconv.ParseInt(string(pos), 10, 64)
	if err != nil {
		return 0
	}
	return i
}

func (s *Storage) SetPos(pos int64) error {
	tx := s.db.NewTransaction(true)
	if err := tx.SetEntry(badger.NewEntry([]byte("pos"), []byte(fmt.Sprint(pos)))); err != nil {
		return xerrors.Errorf("failed to set pos: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return xerrors.Errorf("failed to commit: %w", err)
	}
	return nil
}

func (s *Storage) Node(id int64) ([]byte, error) {
	value, err := s.getItemByKey(s.nodeID(id))
	if err != nil {
		return nil, xerrors.Errorf("failed to get node: %w", err)
	}
	return value, nil
}

func (s *Storage) AddNode(id int64, v []byte) error {
	if err := s.setItem(s.nodeID(id), v); err != nil {
		return xerrors.Errorf("failed to add node: %w", err)
	}
	return nil
}

func (s *Storage) AddNodeInWay(id int64) error {
	if err := s.setItem(s.nodeIDInWay(id), []byte{1}); err != nil {
		return xerrors.Errorf("failed to add node in way: %w", err)
	}
	return nil
}

func (s *Storage) AddWayInMember(id int64) error {
	if err := s.setItem(s.wayIDInMember(id), []byte{1}); err != nil {
		return xerrors.Errorf("failed to add way in member: %w", err)
	}
	return nil
}

func (s *Storage) AddWay(id int64, v []byte) error {
	if err := s.setItem(s.wayID(id), v); err != nil {
		return xerrors.Errorf("failed to add way: %w", err)
	}
	return nil
}

func (s *Storage) AddRelation(id int64, v []byte) error {
	if err := s.setItem(s.relationID(id), v); err != nil {
		return xerrors.Errorf("failed to add relation: %w", err)
	}
	return nil
}

func (s *Storage) Way(id int64) ([]byte, error) {
	value, err := s.getItemByKey(s.wayID(id))
	if err != nil {
		return nil, xerrors.Errorf("failed to get way: %w", err)
	}
	return value, nil
}

func (s *Storage) Relation(id int64) ([]byte, error) {
	value, err := s.getItemByKey(s.relationID(id))
	if err != nil {
		return nil, xerrors.Errorf("failed to get relation: %w", err)
	}
	return value, nil
}

func (s *Storage) ExistsNode(id int64) bool {
	found, _ := s.Node(id)
	return len(found) > 1
}

func (s *Storage) ExistsNodeInWay(id int64) bool {
	found, _ := s.getItemByKey(s.nodeIDInWay(id))
	return len(found) > 0 && found[0] == 1
}

func (s *Storage) ExistsWayInMember(id int64) bool {
	found, _ := s.getItemByKey(s.wayIDInMember(id))
	return len(found) > 0 && found[0] == 1
}

func (s *Storage) ExistsWay(id int64) bool {
	found, _ := s.Way(id)
	return len(found) > 1
}

func (s *Storage) ExistsRelation(id int64) bool {
	found, _ := s.Relation(id)
	return len(found) > 1
}

func (s *Storage) All(f func([]byte, []byte) error) error {
	stream := s.db.NewStream()
	stream.Send = func(kvl *pb.KVList) error {
		kvs := kvl.GetKv()
		for _, kv := range kvs {
			k := kv.GetKey()
			v := kv.GetValue()
			if err := f(k, v); err != nil {
				return xerrors.Errorf("failed to get value: %w", err)
			}
		}
		return nil
	}
	if err := stream.Orchestrate(context.Background()); err != nil {
		return xerrors.Errorf("failed to streaming: %w", err)
	}
	return nil
}

func (s *Storage) AllNodes(f func([]byte) error) error {
	stream := s.db.NewStream()
	stream.Prefix = []byte("node")
	stream.Send = func(kvl *pb.KVList) error {
		kvs := kvl.GetKv()
		for _, kv := range kvs {
			v := kv.GetValue()
			if err := f(v); err != nil {
				return xerrors.Errorf("failed to get value: %w", err)
			}
		}
		return nil
	}
	if err := stream.Orchestrate(context.Background()); err != nil {
		return xerrors.Errorf("failed to streaming: %w", err)
	}
	return nil
}

func (s *Storage) AllWays(f func([]byte) error) error {
	stream := s.db.NewStream()
	stream.Prefix = []byte("way")
	stream.Send = func(kvl *pb.KVList) error {
		for _, kv := range kvl.GetKv() {
			if err := f(kv.GetValue()); err != nil {
				return xerrors.Errorf("failed to get value: %w", err)
			}
		}
		return nil
	}
	if err := stream.Orchestrate(context.Background()); err != nil {
		return xerrors.Errorf("failed to streaming: %w", err)
	}
	return nil
}

func (s *Storage) AllRelations(f func([]byte) error) error {
	stream := s.db.NewStream()
	stream.Prefix = []byte("rel")
	stream.Send = func(kvl *pb.KVList) error {
		for _, kv := range kvl.GetKv() {
			if err := f(kv.GetValue()); err != nil {
				return xerrors.Errorf("failed to get value: %w", err)
			}
		}
		return nil
	}
	if err := stream.Orchestrate(context.Background()); err != nil {
		return xerrors.Errorf("failed to streaming: %w", err)
	}
	return nil
}
