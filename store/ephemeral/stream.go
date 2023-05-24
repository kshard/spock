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

	"github.com/fogfish/golem/trait/pair"
	"github.com/fogfish/golem/trait/seq"
	"github.com/fogfish/skiplist"
	"github.com/kshard/spock"
	"github.com/kshard/xsd"
)

type notSupported struct{ spock.Pattern }

func (err notSupported) Error() string { return fmt.Sprintf("not supported %s", err.Pattern.Dump()) }
func (notSupported) NotSupported()     {}

// TODO: query to stream builder
// - e.g. limit stream
// - skiplist define seq type (like iterator but only value)

type X[C skiplist.Key] interface {
	L1() (uint64, uint64)
	L2() uint64
	ToSPOCK(uint64, seq.Seq[C]) seq.Seq[spock.SPOCK]
}

// helper function to query prefix of the index
func queryIRI[A, B, C skiplist.Key](
	qA *spock.Predicate[A],
	qB *spock.Predicate[B],
	qC *spock.Predicate[C],
	lv X[C],
	list *skiplist.Map[uint64, *skiplist.Set[C]],
) seq.Seq[spock.SPOCK] {
	if qA != nil && qA.Clause == spock.EQ && qB != nil && qB.Clause == spock.EQ {
		ab := lv.L2()
		__x, has := list.Get(ab)
		if !has {
			return nil
		}

		return lv.ToSPOCK(ab, queryXSD(qC, __x))
	}

	if qA != nil && qA.Clause == spock.EQ {
		ab, ab1 := lv.L1()
		a__ := pair.TakeWhile[uint64, *skiplist.Set[C]](
			skiplist.ForMap(list, list.Successors(ab)),
			func(ab uint64, __x *skiplist.Set[C]) bool { return ab < ab1 },
		)
		return pair.JoinSeq(a__,
			func(ab uint64, __x *skiplist.Set[C]) seq.Seq[spock.SPOCK] {
				return lv.ToSPOCK(ab, queryXSD(qC, __x))
			},
		)
	}

	a__ := skiplist.ForMap(list, list.Keys())
	return pair.JoinSeq[uint64, *skiplist.Set[C]](a__,
		func(ab uint64, __x *skiplist.Set[C]) seq.Seq[spock.SPOCK] {
			return lv.ToSPOCK(ab, queryXSD(qC, __x))
		},
	)
}

// 	if qA != nil && qB != nil {
// 		sp := uint64(*qA)<<32 | *qB
// 		__x, has := list.Get(sp)
// 		if !has {
// 			return nil
// 		}

// 		return __x

// 		return toSPO(sp, queryXSD(q.O, __x)), nil
// 	}

// 	var val pair.Seq[uint64, __o]

// 	if q.S != nil && q.S.Clause == spock.EQ {
// 		sp := uint64(q.S.Value) << 32
// 		up := uint64(q.S.Value+1) << 32
// 		val = skiplist.ForMap(store.spo, store.spo.Successors(sp))
// 		val = pair.TakeWhile(val,
// 			func(k uint64, v __o) bool { return k < up },
// 		)
// 	}

// 	// if qA !=

// 	// key := uint64(0)
// 	// if pred != nil && pred.Clause == spock.EQ {
// 	// 	key = uint64(qA.Value) << 32
// 	// }

// 	// seq := skiplist.ForMap(list, list.Successors(key))
// 	// if key != 0 {
// 	// 	max := uint64(pred.Value+1) << 32
// 	// 	seq = skiplist.TakeWhile(seq,
// 	// 		func(k uint64, v C) bool {
// 	// 			return k < max
// 	// 		},
// 	// 	)
// 	// }

// 	// var seq *skiplist.Iterator[s, B]

// 	// switch {
// 	// case pred == nil:
// 	// 	seq = skiplist.Values(list)
// 	// case pred.Clause == spock.EQ:
// 	// 	return NewValueSeq(list, pred.Value).(Seq[A, B])
// 	// default:
// 	// 	panic(fmt.Errorf("xsd.AnyURI do not support %s", pred))
// 	// }

// 	// if seq == nil {
// 	// 	return nil
// 	// }

// 	// return Seq[s, B](seq).(Seq[A, B])
// 	return nil
// }

func queryXSD[A skiplist.Key](
	pred *spock.Predicate[A],
	list *skiplist.Set[A],
) seq.Seq[A] {
	if pred != nil && pred.Clause == spock.EQ {
		if list.Has(pred.Value) {
			return seq.From(pred.Value)
		}
		return nil
	}

	return skiplist.ForSet(list, list.Values())
}

func (store *Store) streamSPO(q spock.Pattern) (seq.Seq[spock.SPOCK], error) {
	return queryIRI[xsd.AnyURI, xsd.AnyURI, xsd.Symbol](
		q.S, q.P, q.O, querySPO(q), store.spo,
	), nil
}

func (store *Store) streamSOP(q spock.Pattern) (seq.Seq[spock.SPOCK], error) {
	return queryIRI[xsd.AnyURI, xsd.Symbol, xsd.AnyURI](
		q.S, q.O, q.P, querySOP(q), store.sop,
	), nil
}

// func (store *Store) streamPSO(q spock.Pattern) (spock.Stream, error) {
// 	return newIterator[p, s, o](queryPSO(q), store.pso), nil
// }

// func (store *Store) streamPOS(q spock.Pattern) (spock.Stream, error) {
// 	return newIterator[p, o, s](queryPOS(q), store.pos), nil
// }

// func (store *Store) streamOSP(q spock.Pattern) (spock.Stream, error) {
// 	return newIterator[o, s, p](queryOSP(q), store.osp), nil
// }

// func (store *Store) streamOPS(q spock.Pattern) (spock.Stream, error) {
// 	return newIterator[o, p, s](queryOPS(q), store.ops), nil
// }

type querySPO spock.Pattern

func (q querySPO) L1() (uint64, uint64) {
	return uint64(q.S.Value) << 32, uint64(q.S.Value+1) << 32
}

func (q querySPO) L2() uint64 {
	return uint64(q.S.Value)<<32 | uint64(q.P.Value)
}

func (q querySPO) ToSPOCK(sp uint64, o seq.Seq[xsd.Symbol]) seq.Seq[spock.SPOCK] {
	return seq.Map(o,
		func(o xsd.Symbol) spock.SPOCK {
			return spock.SPOCK{
				S: xsd.AnyURI(sp >> 32),
				P: xsd.AnyURI(sp & (0xffffffff)),
				O: o,
			}
		},
	)
}

type querySOP spock.Pattern

func (q querySOP) L1() (uint64, uint64) {
	return uint64(q.S.Value) << 32, uint64(q.S.Value+1) << 32
}

func (q querySOP) L2() uint64 {
	return uint64(q.S.Value)<<32 | uint64(q.O.Value)
}

func (q querySOP) ToSPOCK(so uint64, p seq.Seq[xsd.AnyURI]) seq.Seq[spock.SPOCK] {
	return seq.Map(p,
		func(p xsd.AnyURI) spock.SPOCK {
			return spock.SPOCK{
				S: xsd.AnyURI(so >> 32),
				P: p,
				O: xsd.Symbol(so & (0xffffffff)),
			}
		},
	)
}
