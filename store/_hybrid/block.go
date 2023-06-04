package hybrid

import "github.com/fogfish/skiplist"

//
// Memory Block
//

//
// Disk Block
//

// Virtual Memory VMem
// VHeap

type VHeap[K, V any] struct {
	ID   h
	Size int
	Heap *skiplist.SkipList[K, V]
}

//
// [prefix, seq] -> [spo] {96, 96, 96} -> 300 B -> 1 K triples
//
// Page ID -> subject ID // what if page ID spo
//
// 1. Overflow
//   a) new o : establishes new node/page  ... s1 ... s1 ...
//   b) new p
//   c) new s : established new node/page  ... s1 ... s2 ...
//
