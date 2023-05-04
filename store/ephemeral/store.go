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
	"math/rand"
	"time"

	"github.com/fogfish/guid/v2"
	"github.com/fogfish/skiplist"
	"github.com/kshard/spock"
	"github.com/kshard/xsd"
)

// Store is the instance of knowledge storage
type Store struct {
	size   int
	random rand.Source
	spo    spo
	sop    sop
	pso    pso
	pos    pos
	osp    osp
	ops    ops
}

// Create new instance of knowledge storage
func New() *Store {
	rnd := rand.NewSource(time.Now().UnixNano())
	return &Store{
		random: rnd,
		spo:    newSPO(rnd),
		sop:    newSOP(rnd),
		pso:    newPSO(rnd),
		pos:    newPOS(rnd),
		osp:    newOSP(rnd),
		ops:    newOPS(rnd),
	}
}

// Size returns number of knowledge statements in the store
func Size(store *Store) int {
	return store.size
}

func Add(store *Store, bag spock.Bag) {
	for _, spock := range bag {
		Put(store, spock)
	}
}

func Put(store *Store, spock spock.SPOCK) {
	spock.K = guid.L(guid.Clock)

	_po, _op := ensureForS(store, spock.S)
	_so, _os := ensureForP(store, spock.P)
	_sp, _ps := ensureForO(store, spock.O)

	putO(store, _po, _so, spock)
	putP(store, _op, _sp, spock)
	putS(store, _os, _ps, spock)

	store.size++
}

func ensureForS(store *Store, s s) (_po, _op) {
	_po, has := skiplist.Lookup(store.spo, s)
	if !has {
		_po = newPO(store.random)
		skiplist.Put(store.spo, s, _po)
	}

	_op, has := skiplist.Lookup(store.sop, s)
	if !has {
		_op = newOP(store.random)
		skiplist.Put(store.sop, s, _op)
	}
	return _po, _op
}

func ensureForP(store *Store, p p) (_so, _os) {
	_so, has := skiplist.Lookup(store.pso, p)
	if !has {
		_so = newSO(store.random)
		skiplist.Put(store.pso, p, _so)
	}

	_os, has := skiplist.Lookup(store.pos, p)
	if !has {
		_os = newOS(store.random)
		skiplist.Put(store.pos, p, _os)
	}
	return _so, _os
}

func ensureForO(store *Store, o o) (_sp, _ps) {
	_sp, has := skiplist.Lookup(store.osp, o)
	if !has {
		_sp = newSP(store.random)
		skiplist.Put(store.osp, o, _sp)
	}

	_ps, has := skiplist.Lookup(store.ops, o)
	if !has {
		_ps = newPS(store.random)
		skiplist.Put(store.ops, o, _ps)
	}
	return _sp, _ps
}

func putO(store *Store, _po _po, _so _so, spock spock.SPOCK) {
	__o, has := skiplist.Lookup(_po, spock.P)
	if !has {
		__o = newO(store.random)
		skiplist.Put(_po, spock.P, __o)
		skiplist.Put(_so, spock.S, __o)
	}

	skiplist.Put(__o, spock.O, struct{}{}) // spock.K)
}

func putP(store *Store, _op _op, _sp _sp, spock spock.SPOCK) {
	__p, has := skiplist.Lookup(_sp, spock.S)
	if !has {
		__p = newP(store.random)
		skiplist.Put(_op, spock.O, __p)
		skiplist.Put(_sp, spock.S, __p)
	}

	skiplist.Put(__p, spock.P, struct{}{}) // spock.K)
}

func putS(store *Store, _os _os, _ps _ps, spock spock.SPOCK) {
	__s, has := skiplist.Lookup(_ps, spock.P)
	if !has {
		__s = newS(store.random)
		skiplist.Put(_os, spock.O, __s)
		skiplist.Put(_ps, spock.P, __s)
	}

	skiplist.Put(__s, spock.S, struct{}{}) // spock.K)
}

func Match(store *Store, q spock.Pattern) (spock.Stream, error) {
	if q.HintForS != spock.HINT_MATCH && q.HintForS != spock.HINT_NONE {
		return nil, &notSupported{q}
	}

	if q.HintForP != spock.HINT_MATCH && q.HintForP != spock.HINT_NONE {
		return nil, &notSupported{q}
	}

	if q.HintForO != spock.HINT_MATCH && q.HintForO != spock.HINT_NONE && q.O.Value.XSDType() == xsd.XSD_ANYURI {
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
