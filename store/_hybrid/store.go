package hybrid

import (
	"math/rand"
	"time"

	"github.com/fogfish/skiplist"
	"github.com/fogfish/skiplist/ord"
	"github.com/kshard/spock"
	"github.com/kshard/xsd"
)

// Store is the instance of knowledge storage
type Store struct {
	size   int
	random rand.Source
	heap   *Heap
	spo    hspo
}

// Create new instance of knowledge storage
func New() *Store {
	rnd := rand.NewSource(time.Now().UnixNano())
	return &Store{
		random: rnd,
		heap:   &Heap{},
		spo:    skiplist.New[h, spo](ord.UInt32, rnd),
	}
}

func Config(store *Store) {
	skiplist.Put(store.spo, 0, nil)
}

func Put(store *Store, spock spock.SPOCK) {
	_po /*, _op*/ := ensureForS(store, spock.S)

	putO(store, _po, spock)
}

func ensureForS(store *Store, x s) _po { // (_po, _op) {
	h := uint32(0) // h <- x

	spo, has := skiplist.Lookup(store.spo, h)
	if !has || spo == nil {
		spo = Swap(store, h)
		skiplist.Put(store.spo, h, spo)
	}

	_po, has := skiplist.Lookup(spo, x)
	if !has {
		_po = skiplist.New[p, __o](xsd.OrdAnyURI, store.random)
		skiplist.Put(spo, x, _po)
	}

	// _op, has := skiplist.Lookup(store.sop, s)
	// if !has {
	// 	_op = newOP(store.random)
	// 	skiplist.Put(store.sop, s, _op)
	// }
	return _po //, _op
}

func putO(store *Store, _po _po, spock spock.SPOCK) {
	__o, has := skiplist.Lookup(_po, spock.P)
	if !has {
		__o = skiplist.New[o, k](xsd.OrdValue, store.random)
		skiplist.Put(_po, spock.P, __o)
		// skiplist.Put(_so, spock.S, __o)
	}

	skiplist.Put(__o, spock.O, struct{}{}) // spock.K)
}

func X(store *Store) (spock.Stream, error) {
	q := querySPO{
		store:   store,
		Pattern: spock.Pattern{},
	}

	return newIterator[s, p, o](q, store.spo), nil
}
