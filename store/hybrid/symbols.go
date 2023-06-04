package hybrid

import "github.com/fogfish/segment"

type Symbol uint32

// Symbol is a global dictionary of literals
type symbol struct {
	kv *segment.Map[uint32, uint32]
}

func (symbol *symbol) ToSymbol() {}

// bucket (id / bucketSize) []seq

// trade off writes for read
// read from global map[uint32]string but read miss should cause load
// write (hash) string -> segment

//
//
