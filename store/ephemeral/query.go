/*

  Knowledge Graph: SPOCK
  Copyright (C) 2016 - 2023 Dmitry Kolesnikov

  This program is free software: you can redistribute it and/or modify
  it under the terms of the GNU Affero General Public License as published
  by the Free Software Foundation, either version 3 of the License, or
  (at your option) any later version.

  This program is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU Affero General Public License for more details.

  You should have received a copy of the GNU Affero General Public License
  along with this program.  If not, see <https://www.gnu.org/licenses/>.

*/

package ephemeral

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
	case pred.Clause == spock.PQ:
		_, after := skiplist.Split(list, pred.Value)
		if after == nil {
			return nil
		}
		return NewTakeWhile[xsd.Value, B](
			func(x xsd.Value) bool { return xsd.HasPrefix(x, pred.Value) },
			after,
		).(Seq[A, B])
	case pred.Clause == spock.IN:
		seq = skiplist.Range(list, pred.Value, pred.Other)
	case pred.Clause == spock.LT:
		before, _ := skiplist.Split(list, pred.Value)
		if before == nil {
			return nil
		}
		return NewDropWhileType[B](pred.Value.XSDType(), before).(Seq[A, B])
	case pred.Clause == spock.GT:
		_, after := skiplist.Split(list, pred.Value)
		if after == nil {
			return nil
		}
		return NewTakeWhileType[B](pred.Value.XSDType(), after).(Seq[A, B])
	}

	if seq == nil {
		return nil
	}

	return Seq[xsd.Value, B](seq).(Seq[A, B])
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

type takeWhile[A, B any] struct {
	Seq[A, B]
	f func(A) bool
}

func NewTakeWhile[A, B any](f func(A) bool, seq Seq[A, B]) Seq[A, B] {
	return &takeWhile[A, B]{Seq: seq, f: f}
}

func (seq *takeWhile[A, B]) Next() bool {
	if !seq.Seq.Next() {
		return false
	}

	if key, _ := seq.Seq.Head(); !seq.f(key) {
		return false
	}

	return true
}

// take sequence elements while xsd.Value belongs to same category (type)
type takeWhileType[T any] struct {
	Seq[xsd.Value, T]
	cat xsd.Symbol
}

func NewTakeWhileType[T any](cat xsd.Symbol, seq Seq[xsd.Value, T]) Seq[xsd.Value, T] {
	return &takeWhileType[T]{Seq: seq, cat: cat}
}

func (seq *takeWhileType[T]) Next() bool {
	if !seq.Seq.Next() {
		return false
	}

	if key, _ := seq.Seq.Head(); key.XSDType() != seq.cat {
		return false
	}

	return true
}

type dropWhileType[T any] struct {
	Seq[xsd.Value, T]
	cat xsd.Symbol
}

func NewDropWhileType[T any](cat xsd.Symbol, seq Seq[xsd.Value, T]) Seq[xsd.Value, T] {
	return &dropWhileType[T]{Seq: seq, cat: cat}
}

func (seq *dropWhileType[T]) Next() bool {
	for {
		if !seq.Seq.Next() {
			return false
		}

		if key, _ := seq.Seq.Head(); key.XSDType() == seq.cat {
			return true
		}
	}
}

// executes query against ⟨s, p, o⟩ data structure
type querySPO spock.Pattern

func (q querySPO) L1(list *skiplist.SkipList[s, _po]) Seq[s, _po] {
	return queryIRI[s](q.S, list)
}

func (q querySPO) L2(list *skiplist.SkipList[p, __o]) Seq[p, __o] {
	return queryIRI[p](q.P, list)
}

func (q querySPO) L3(list *skiplist.SkipList[o, k]) Seq[o, k] {
	return queryXSD[o](q.O, list)
}

func (q querySPO) ToSPOCK(s s, p p, o o) spock.SPOCK {
	return spock.SPOCK{S: s, P: p, O: o}
}

// executes query against ⟨s, o, p⟩ data structure
type querySOP spock.Pattern

func (q querySOP) L1(list *skiplist.SkipList[s, _op]) Seq[s, _op] {
	return queryIRI[s](q.S, list)
}

func (q querySOP) L2(list *skiplist.SkipList[o, __p]) Seq[o, __p] {
	return queryXSD[o](q.O, list)
}

func (q querySOP) L3(list *skiplist.SkipList[p, k]) Seq[p, k] {
	return queryIRI[p](q.P, list)
}

func (q querySOP) ToSPOCK(s s, o o, p p) spock.SPOCK {
	return spock.SPOCK{S: s, P: p, O: o}
}

// executes query against ⟨p, s, o⟩ data structure
type queryPSO spock.Pattern

func (q queryPSO) L1(list *skiplist.SkipList[p, _so]) Seq[p, _so] {
	return queryIRI[p](q.P, list)
}

func (q queryPSO) L2(list *skiplist.SkipList[s, __o]) Seq[s, __o] {
	return queryIRI[s](q.S, list)
}

func (q queryPSO) L3(list *skiplist.SkipList[o, k]) Seq[o, k] {
	return queryXSD[o](q.O, list)
}

func (q queryPSO) ToSPOCK(p p, s s, o o) spock.SPOCK {
	return spock.SPOCK{S: s, P: p, O: o}
}

// executes query against ⟨p, o, s⟩ data structure
type queryPOS spock.Pattern

func (q queryPOS) L1(list *skiplist.SkipList[p, _os]) Seq[p, _os] {
	return queryIRI[p](q.P, list)
}

func (q queryPOS) L2(list *skiplist.SkipList[o, __p]) Seq[o, __p] {
	return queryXSD[o](q.O, list)
}

func (q queryPOS) L3(list *skiplist.SkipList[s, k]) Seq[s, k] {
	return queryIRI[s](q.S, list)
}

func (q queryPOS) ToSPOCK(p p, o o, s s) spock.SPOCK {
	return spock.SPOCK{S: s, P: p, O: o}
}

// executes query against ⟨o, p, s⟩ data structure
type queryOPS spock.Pattern

func (q queryOPS) L1(list *skiplist.SkipList[o, _ps]) Seq[o, _ps] {
	return queryXSD[o](q.O, list)
}

func (q queryOPS) L2(list *skiplist.SkipList[p, __s]) Seq[p, __s] {
	return queryIRI[p](q.P, list)
}

func (q queryOPS) L3(list *skiplist.SkipList[s, k]) Seq[s, k] {
	return queryIRI[s](q.S, list)
}

func (q queryOPS) ToSPOCK(o o, p p, s s) spock.SPOCK {
	return spock.SPOCK{S: s, P: p, O: o}
}

// executes query against ⟨o, s, p⟩ data structure
type queryOSP spock.Pattern

func (q queryOSP) L1(list *skiplist.SkipList[o, _ps]) Seq[o, _ps] {
	return queryXSD[o](q.O, list)
}

func (q queryOSP) L2(list *skiplist.SkipList[s, __p]) Seq[s, __p] {
	return queryIRI[s](q.S, list)
}

func (q queryOSP) L3(list *skiplist.SkipList[p, k]) Seq[p, k] {
	return queryIRI[p](q.P, list)
}

func (q queryOSP) ToSPOCK(o o, s s, p p) spock.SPOCK {
	return spock.SPOCK{S: s, P: p, O: o}
}
