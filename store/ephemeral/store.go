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
	spo  *skiplist.Map[sp, __o]
	sop  *skiplist.Map[so, __p]
	pso  *skiplist.Map[ps, __o]
	pos  *skiplist.Map[po, __s]
	osp  *skiplist.Map[os, __p]
	ops  *skiplist.Map[op, __s]
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

func (store *Store) Put(spock spock.SPOCK) error {
	spock.K = guid.L(guid.Clock)

	kind := spock.O.XSDType()
	if kind != xsd.XSD_ANYURI /*|| kind != xsd.XSD_SYMBOL*/ {
		return fmt.Errorf("not supported %T", spock.O)
	}

	o, ok := spock.O.(xsd.AnyURI)
	if !ok {
		return fmt.Errorf("not supported %T", spock.O)
	}

	store.putSuffixO(spock.S, spock.P, o)
	store.putSuffixS(spock.S, spock.P, o)
	store.putSuffixP(spock.S, spock.P, o)

	store.size++
	return nil
}

func (store *Store) putSuffixO(s, p, o xsd.AnyURI) bool {
	sp := uint64(s)<<32 | uint64(p)
	ps := uint64(p)<<32 | uint64(s)
	__o, node := store.spo.Get(sp)
	if node == nil {
		__o = skiplist.NewSet[xsd.AnyURI]()
		store.spo.Put(sp, __o)
		store.pso.Put(ps, __o)
	}

	has, _ := __o.Add(o)
	return has
}

func (store *Store) putSuffixS(s, p, o xsd.AnyURI) bool {
	po := uint64(p)<<32 | uint64(o)
	op := uint64(o)<<32 | uint64(p)
	__s, node := store.pos.Get(po)
	if node == nil {
		__s = skiplist.NewSet[xsd.AnyURI]()
		store.pos.Put(po, __s)
		store.ops.Put(op, __s)
	}

	has, _ := __s.Add(s)
	return has
}

func (store *Store) putSuffixP(s, p, o xsd.AnyURI) bool {
	so := uint64(s)<<32 | uint64(o)
	os := uint64(o)<<32 | uint64(s)
	__p, node := store.sop.Get(so)
	if node == nil {
		__p = skiplist.NewSet[xsd.AnyURI]()
		store.sop.Put(so, __p)
		store.osp.Put(os, __p)
	}

	has, _ := __p.Add(p)
	return has
}

func (store *Store) Match(q spock.Pattern[xsd.AnyURI]) (seq.Seq[spock.SPOCK], error) {
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
	case spock.STRATEGY_PSO:
		return store.streamPSO(q)
	case spock.STRATEGY_POS:
		return store.streamPOS(q)
	case spock.STRATEGY_OSP:
		return store.streamOSP(q)
	case spock.STRATEGY_OPS:
		return store.streamOPS(q)
	default:
		panic(fmt.Errorf("unknown strategy"))
	}
}
