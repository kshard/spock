package symbol

import (
	"context"
	"crypto/sha256"
	"fmt"
	"strings"
	"sync"

	"github.com/fogfish/curie"
	"github.com/fogfish/guid/v2"
	"github.com/kshard/xsd"
)

type Store interface {
	Match(context.Context, *Symbols, ...interface{ MatchOpt() }) ([]*Symbols, interface{ MatchOpt() }, error)
	Put(context.Context, *Symbols, ...interface{ ConditionExpression(*Symbols) }) error
}

type Namespace struct {
	sync.Mutex

	id    curie.IRI
	clock guid.Chronos

	store     Store
	syncAfter guid.K

	bySymbol map[Symbol]xsd.AnyURI
	byAnyURI map[xsd.AnyURI]Symbol
}

func New(id curie.IRI, store Store) *Namespace {
	h := sha256.New()
	h.Write([]byte(id))
	hash := h.Sum(nil)
	node := uint64(hash[0])<<24 | uint64(hash[1])<<16 | uint64(hash[2])<<8 | uint64(hash[3])
	clock := guid.NewClock(guid.WithNodeID(node))

	return &Namespace{
		id:        id,
		clock:     clock,
		store:     store,
		syncAfter: guid.G(clock),
		bySymbol:  map[Symbol]xsd.AnyURI{},
		byAnyURI:  map[xsd.AnyURI]Symbol{},
	}
}

func (ns *Namespace) ToSymbol(uri xsd.AnyURI) Symbol {
	ns.Lock()
	defer ns.Unlock()

	symb, has := ns.byAnyURI[uri]
	if !has {
		symb = Symbol(guid.G(ns.clock))
		ns.bySymbol[symb] = uri
		ns.byAnyURI[uri] = symb
	}

	return symb
}

func (ns *Namespace) ToAnyURI(sym Symbol) xsd.AnyURI {
	ns.Lock()
	uri, has := ns.bySymbol[sym]
	ns.Unlock()

	if has {
		return uri
	}

	// NiiSbqR8dWfVuJ7S
	// NiiSbqR8dWfVuJ7S
	// fmt.Println("==>> fuq " + sym.String())

	return xsd.AnyURI(xsd.XSD_NIL)
}

func (ns *Namespace) FromString(sym string) xsd.AnyURI {
	k, err := guid.FromStringG(sym)
	if err != nil {
		return xsd.AnyURI(xsd.XSD_NIL)
	}

	return ns.ToAnyURI(Symbol(k))
}

// TODO:
//  - lazy loading base on consistent hashing
//  - this approach causes "overheat" of dynamo
//  - all symbols are loaded at once

func (ns *Namespace) Read(ctx context.Context) error {
	ns.Lock()
	defer ns.Unlock()

	var cursor interface{ MatchOpt() }
	key := &Symbols{Namespace: ns.id}

	for ok := true; ok; ok = cursor != nil {
		seq, cur, err := ns.store.Match(ctx, key, cursor)
		if err != nil {
			return err
		}

		err = ns.decodeSymbols(seq)
		if err != nil {
			return err
		}

		cursor = cur
	}

	return nil
}

func (ns *Namespace) decodeSymbols(seq []*Symbols) error {
	for _, page := range seq {
		for _, symbol := range page.Symbols {
			pair := strings.Split(symbol, "|")
			// if pair[1] == "ub:GraduateStudent" {
			// fmt.Println("===>>>> " + symbol)
			// }

			sym, _ := guid.FromStringG(pair[0])
			uri := xsd.ToAnyURI(curie.IRI(pair[1]))
			ns.bySymbol[Symbol(sym)] = uri
			ns.byAnyURI[uri] = Symbol(sym)
		}
	}

	return nil
}

func (ns *Namespace) Write(ctx context.Context) error {
	ns.Lock()
	defer ns.Unlock()

	fmt.Println("==> writing")
	writer := func(seq []string) error {
		return ns.store.Put(ctx, &Symbols{
			Namespace: ns.id,
			ID:        curie.IRI(guid.G(ns.clock).String()),
			Symbols:   seq,
		})
	}

	ssz := 0
	seq := []string{}
	for symbol, anyURI := range ns.bySymbol {
		k := guid.K(symbol)
		if guid.After(k, ns.syncAfter) {
			fmt.Printf("==> %s\n", anyURI.String())
			val := k.String() + "|" + anyURI.String()
			ssz = ssz + len(val)
			seq = append(seq, val)

			if ssz > 300*1024 {
				if err := writer(seq); err != nil {
					return err
				}
				seq = []string{}
				ssz = 0
			}
		} else {
			fmt.Printf("==>>\nk = %s (%b, %b)\n, a = %s (%b, %b)\n", k.String(), k.Hi, k.Lo, ns.syncAfter, ns.syncAfter.Hi, ns.syncAfter.Lo)
		}
	}

	if len(seq) != 0 {
		if err := writer(seq); err != nil {
			return err
		}
	}

	fmt.Println("==> written")

	ns.syncAfter = guid.G(ns.clock)
	return nil
}
