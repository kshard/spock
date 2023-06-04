package hybrid

import (
	"fmt"

	"github.com/fogfish/skiplist"
	"github.com/kshard/xsd"
)

type Heap struct {
	// page map[h]int
}

func NewHeap() *Heap {
	// TODO: load from store page index
	return &Heap{
		// page: map[h]int{},
	}
}

func Swap(store *Store, h h) *skiplist.SkipList[s, _po] {
	fmt.Printf("==> new page %v\n", h)
	// TODO: page miss
	return skiplist.New[s, _po](xsd.OrdAnyURI, store.random)
}

// page ID: X byte + # Seq
