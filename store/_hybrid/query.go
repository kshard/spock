package hybrid

import (
	"fmt"

	"github.com/fogfish/skiplist"
	"github.com/kshard/spock"
	"github.com/kshard/xsd"
)

// Each query results with sequence of "elements".
// This interface defines generic sequence, abstracting skiplist.Iterator
type Seq[K, V any] interface {
	Head() (K, V)
	Next() bool
}

// helper function to query the skiplist where key is curie.IRI
func queryShard[A, B any](
	pred *spock.Predicate[h],
	list *skiplist.SkipList[h, B],
) Seq[A, B] {
	var seq *skiplist.Iterator[h, B]

	switch {
	case pred == nil:
		seq = skiplist.Values(list)
	case pred.Clause == spock.EQ:
		return NewValueSeq(list, pred.Value).(Seq[A, B])
	default:
		panic(fmt.Errorf("xsd.AnyURI do not support %s", pred))
	}

	if seq == nil {
		return nil
	}

	return Seq[h, B](seq).(Seq[A, B])
}

type valueSeq[K, V any] struct {
	key K
	val V
	seq *skiplist.SkipList[K, V]
}

func NewValueSeq[K, V any](seq *skiplist.SkipList[K, V], key K) Seq[K, V] {
	return &valueSeq[K, V]{
		key: key,
		seq: seq,
	}
}

func (seq *valueSeq[K, V]) Head() (K, V) {
	return seq.key, seq.val
}

func (seq *valueSeq[K, V]) Next() bool {
	if seq.seq != nil {
		val, has := skiplist.Lookup(seq.seq, seq.key)
		seq.val = val
		seq.seq = nil
		return has
	}

	return false
}

// helper function to query the skiplist where key is curie.IRI
func queryIRI[A, B any](
	pred *spock.Predicate[s],
	list *skiplist.SkipList[s, B],
) Seq[A, B] {
	var seq *skiplist.Iterator[s, B]

	switch {
	case pred == nil:
		seq = skiplist.Values(list)
	case pred.Clause == spock.EQ:
		return NewValueSeq(list, pred.Value).(Seq[A, B])
	default:
		panic(fmt.Errorf("xsd.AnyURI do not support %s", pred))
	}

	if seq == nil {
		return nil
	}

	return Seq[s, B](seq).(Seq[A, B])
}

// helper function to query the skiplist where key is xsd.Value
func queryXSD[A, B any](
	pred *spock.Predicate[o],
	list *skiplist.SkipList[o, B],
) Seq[A, B] {
	var seq *skiplist.Iterator[o, B]

	switch {
	case pred == nil:
		seq = skiplist.Values(list)
	case pred.Clause == spock.EQ:
		seq = skiplist.Slice(list, pred.Value, 1)
		// case pred.Clause == spock.PQ:
		// 	_, after := skiplist.Split(list, pred.Value)
		// 	if after == nil {
		// 		return nil
		// 	}
		// 	return NewTakeWhile[xsd.Value, B](
		// 		func(x xsd.Value) bool { return xsd.HasPrefix(x, pred.Value) },
		// 		after,
		// 	).(Seq[A, B])
		// case pred.Clause == spock.IN:
		// 	seq = skiplist.Range(list, pred.Value, pred.Other)
		// case pred.Clause == spock.LT:
		// 	before, _ := skiplist.Split(list, pred.Value)
		// 	if before == nil {
		// 		return nil
		// 	}
		// 	return NewDropWhileType[B](pred.Value.XSDType(), before).(Seq[A, B])
		// case pred.Clause == spock.GT:
		// 	_, after := skiplist.Split(list, pred.Value)
		// 	if after == nil {
		// 		return nil
		// 	}
		// 	return NewTakeWhileType[B](pred.Value.XSDType(), after).(Seq[A, B])
	}

	if seq == nil {
		return nil
	}

	return Seq[xsd.Value, B](seq).(Seq[A, B])
}

// executes query against ⟨s, p, o⟩ data structure
type querySPO struct {
	store *Store
	spock.Pattern
}

func (q querySPO) L0(list *skiplist.SkipList[h, spo]) Seq[h, spo] {
	return queryShard[h](nil, list)
}

func (q querySPO) L1(list *skiplist.SkipList[s, _po]) Seq[s, _po] {
	return queryIRI[s](q.S, list)
}

func (q querySPO) L2(list *skiplist.SkipList[p, __o]) Seq[p, __o] {
	return queryIRI[p](q.P, list)
}

func (q querySPO) L3(list *skiplist.SkipList[o, k]) Seq[o, k] {
	return queryXSD[o](q.O, list)
}

func (q querySPO) Swap(h h) *skiplist.SkipList[s, _po] {
	return Swap(q.store, h)
}

func (q querySPO) ToSPOCK(s s, p p, o o) spock.SPOCK {
	return spock.SPOCK{S: s, P: p, O: o}
}
