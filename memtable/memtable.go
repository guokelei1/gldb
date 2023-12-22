package memtable

import (
	"gldb/internal"
	"gldb/skiplist"
)

type MemTable struct {
	table       *skiplist.SkipList
	memoryUsage uint64
}

func NewMemTable() *MemTable {
	return &MemTable{
		table:       skiplist.New(internal.InternalKeyComparator),
		memoryUsage: 0,
	}
}

func (mt *MemTable) NewIterator() *Iterator {
	return &Iterator{
		listIter: mt.table.NewIterator(),
	}
}

func (mt *MemTable) Add(seq uint64, valueType internal.ValueType, key, value []byte) {
	internalKey := internal.NewInternalKey(seq, valueType, key, value)
	mt.table.Insert(internalKey)
	mt.memoryUsage += uint64(16 + len(key) + len(value))
}

func (mt *MemTable) Get(key []byte) ([]byte, error) {
	lookupKey := internal.LookupKey(key)
	it := mt.table.NewIterator()
	it.Seek(lookupKey)
	if it.Valid() {
		internalKey := it.Key().(*internal.InternalKey)
		if internal.UserKeyComparator(key, internalKey.UserKey) == 0 {
			if internalKey.Type == internal.TypeValue {
				return internalKey.UserValue, nil
			} else {
				return nil, internal.ErrDeletion
			}
		}

	}
	return nil, internal.ErrNotFound
}
func (mt *MemTable) ApproximateMemoryUsage() uint64 {
	return mt.memoryUsage
}
