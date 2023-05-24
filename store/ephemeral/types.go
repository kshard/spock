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
	"github.com/fogfish/skiplist"
	"github.com/kshard/xsd"
)

// components of <s,p,o,c,k> triple
// type s = xsd.AnyURI // subject
// type p = xsd.AnyURI // predicate
// type o = xsd.Symbol // object
// type c = float64 // TODO: credibility
// type k = struct{} // TODO: k-order guid.K

type (
	s  = uint32
	p  = uint32
	o  = uint32
	sp = uint64
	so = uint64
	ps = uint64
	po = uint64
	os = uint64
	op = uint64
)

// index types for 3rd faction
type __s = *skiplist.Set[xsd.AnyURI]
type __p = *skiplist.Set[xsd.AnyURI]
type __o = *skiplist.Set[xsd.Symbol]

// index types for 2nd faction
// type _po = *skiplist.SkipList[p, __o]
// type _op = *skiplist.SkipList[o, __p]
// type _so = *skiplist.SkipList[s, __o]
// type _os = *skiplist.SkipList[o, __s]
// type _sp = *skiplist.SkipList[s, __p]
// type _ps = *skiplist.SkipList[p, __s]

// triple indexes
type spo = *skiplist.Map[sp, __o]
type sop = *skiplist.Map[so, __p]
type pso = *skiplist.Map[ps, __o]
type pos = *skiplist.Map[po, __s]
type osp = *skiplist.Map[os, __p]
type ops = *skiplist.Map[op, __s]

// allocators for indexes
// func newS(rnd rand.Source) __s { return skiplist.New[s, k](xsd.OrdAnyURI, rnd) }
// func newP(rnd rand.Source) __p { return skiplist.New[p, k](xsd.OrdAnyURI, rnd) }
// func newO(rnd rand.Source) __o { return skiplist.New[o, k](xsd.OrdValue, rnd) }

// func newPO(rnd rand.Source) _po { return skiplist.New[p, __o](xsd.OrdAnyURI, rnd) }
// func newOP(rnd rand.Source) _op { return skiplist.New[o, __p](xsd.OrdValue, rnd) }
// func newSO(rnd rand.Source) _so { return skiplist.New[s, __o](xsd.OrdAnyURI, rnd) }
// func newOS(rnd rand.Source) _os { return skiplist.New[o, __s](xsd.OrdValue, rnd) }
// func newSP(rnd rand.Source) _sp { return skiplist.New[s, __p](xsd.OrdAnyURI, rnd) }
// func newPS(rnd rand.Source) _ps { return skiplist.New[p, __s](xsd.OrdAnyURI, rnd) }

// func newSPO(rnd rand.Source) spo { return skiplist.New[s, _po](xsd.OrdAnyURI, rnd) }
// func newSOP(rnd rand.Source) sop { return skiplist.New[s, _op](xsd.OrdAnyURI, rnd) }
// func newPSO(rnd rand.Source) pso { return skiplist.New[p, _so](xsd.OrdAnyURI, rnd) }
// func newPOS(rnd rand.Source) pos { return skiplist.New[p, _os](xsd.OrdAnyURI, rnd) }
// func newOSP(rnd rand.Source) osp { return skiplist.New[o, _sp](xsd.OrdValue, rnd) }
// func newOPS(rnd rand.Source) ops { return skiplist.New[o, _ps](xsd.OrdValue, rnd) }
