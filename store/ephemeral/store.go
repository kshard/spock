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

	"github.com/fogfish/golem/trait/seq"
	"github.com/fogfish/guid/v2"
	"github.com/fogfish/skiplist"
	"github.com/kshard/spock"
	"github.com/kshard/xsd"
)

// Store is the instance of knowledge storage
type Store struct {
	size int
	spo  spo
	sop  sop
	pso  pso
	pos  pos
	osp  osp
	ops  ops
}

// Create new instance of knowledge storage
func New() *Store {
	return &Store{
		spo: skiplist.NewMap[sp, __o](),
		sop: skiplist.NewMap[so, __p](),
		pso: skiplist.NewMap[ps, __o](),
		pos: skiplist.NewMap[po, __s](),
		osp: skiplist.NewMap[os, __p](),
		ops: skiplist.NewMap[op, __s](),
	}
}

// Size returns number of knowledge statements in the store
func (store *Store) Size() int {
	return store.size
}

func (store *Store) Add(bag spock.Bag) {
	for _, spock := range bag {
		store.Put(spock)
	}
}

func (store *Store) Put(spock spock.SPOCK) {
	spock.K = guid.L(guid.Clock)

	store.putSuffixO(spock)
	store.putSuffixS(spock)
	store.putSuffixP(spock)

	store.size++
}

func (store *Store) putSuffixO(spock spock.SPOCK) bool {
	sp := uint64(spock.S)<<32 | uint64(spock.P)
	ps := uint64(spock.P)<<32 | uint64(spock.S)
	__o, has := store.spo.Get(sp)
	if !has {
		__o = skiplist.NewSet[xsd.Symbol]()
		store.spo.Put(sp, __o)
		store.pso.Put(ps, __o)
	}

	return __o.Add(spock.O)
}

func (store *Store) putSuffixS(spock spock.SPOCK) bool {
	po := uint64(spock.P)<<32 | uint64(spock.O)
	op := uint64(spock.O)<<32 | uint64(spock.P)
	__s, has := store.pos.Get(po)
	if !has {
		__s = skiplist.NewSet[xsd.AnyURI]()
		store.pos.Put(po, __s)
		store.ops.Put(op, __s)
	}

	return __s.Add(spock.S)
}

func (store *Store) putSuffixP(spock spock.SPOCK) bool {
	so := uint64(spock.S)<<32 | uint64(spock.O)
	os := uint64(spock.O)<<32 | uint64(spock.S)
	__p, has := store.sop.Get(so)
	if !has {
		__p = skiplist.NewSet[xsd.AnyURI]()
		store.sop.Put(so, __p)
		store.osp.Put(os, __p)
	}

	return __p.Add(spock.P)
}

func (store *Store) Match(q spock.Pattern) (seq.Seq[spock.SPOCK], error) {
	if q.HintForS != spock.HINT_MATCH && q.HintForS != spock.HINT_NONE {
		return nil, &notSupported{q}
	}

	if q.HintForP != spock.HINT_MATCH && q.HintForP != spock.HINT_NONE {
		return nil, &notSupported{q}
	}

	if q.HintForO != spock.HINT_MATCH && q.HintForO != spock.HINT_NONE { //&& q.O.Value.XSDType() == xsd.XSD_ANYURI {
		return nil, &notSupported{q}
	}

	switch q.Strategy {
	case spock.STRATEGY_SPO:
		return store.streamSPO(q)
	case spock.STRATEGY_SOP:
		return store.streamSOP(q)
	// case spock.STRATEGY_PSO:
	// 	return store.streamPSO(q)
	// case spock.STRATEGY_POS:
	// 	return store.streamPOS(q)
	// case spock.STRATEGY_OSP:
	// 	return store.streamOSP(q)
	// case spock.STRATEGY_OPS:
	// 	return store.streamOPS(q)
	default:
		panic(fmt.Errorf("unknown strategy"))
	}
}
