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

package dynamo

import (
	"context"

	"github.com/fogfish/curie"
	"github.com/fogfish/dynamo/v2/service/ddb"
	"github.com/kshard/spock"
)

type Writer interface {
	Put(ctx context.Context, store *Store) error
	Cut(ctx context.Context, store *Store) error
}

//
// ⟨ Subject, Predicate, Object ⟩
//

type spo struct {
	G  curie.IRI `dynamodbav:"prefix"`
	SP string    `dynamodbav:"suffix"`
	O  []string  `dynamodbav:"o,stringset"`
}

func (spo spo) HashKey() curie.IRI { return spo.G }
func (spo spo) SortKey() curie.IRI { return curie.IRI(spo.SP) }

func (spo spo) ToSPOCK(symbols Symbols) []spock.SPOCK {
	return decodeSPO(symbols, spo)
}

func (spo spo) Put(ctx context.Context, store *Store) error {
	_, err := store.spo.UpdateWith(ctx,
		ddb.Updater(spo, _spo.Union(spo.O)),
	)
	return err
}

func (spo spo) Cut(ctx context.Context, store *Store) error {
	_, err := store.spo.UpdateWith(ctx,
		ddb.Updater(spo, _spo.Minus(spo.O)),
	)
	return err
}

var (
	_spo = ddb.UpdateFor[spo, []string]()
)

func encodeSPO(symbols Symbols, g curie.IRI, spock spock.SPOCK) spo {
	return spo{
		G:  "sp|" + g,
		SP: encodeII(symbols, spock.S, spock.P),
		O:  []string{encodeValue(symbols, spock.O)},
	}
}

func decodeSPO(symbols Symbols, spo spo) []spock.SPOCK {
	seq := make([]spock.SPOCK, len(spo.O))
	s, p := decodeII(symbols, spo.SP)

	for i, o := range spo.O {
		seq[i].S, seq[i].P, seq[i].O = s, p, decodeValue(symbols, o)
	}

	return seq
}

//
// ⟨ Subject, Object, Predicate ⟩
//

type sop struct {
	G  curie.IRI `dynamodbav:"prefix"`
	SO string    `dynamodbav:"suffix"`
	P  []string  `dynamodbav:"p,stringset"`
}

func (sop sop) HashKey() curie.IRI { return sop.G }
func (sop sop) SortKey() curie.IRI { return curie.IRI(sop.SO) }

func (sop sop) ToSPOCK(symbols Symbols) []spock.SPOCK {
	return decodeSOP(symbols, sop)
}

func (sop sop) Put(ctx context.Context, store *Store) error {
	_, err := store.sop.UpdateWith(ctx,
		ddb.Updater(sop, _sop.Union(sop.P)),
	)
	return err
}

func (sop sop) Cut(ctx context.Context, store *Store) error {
	_, err := store.sop.UpdateWith(ctx,
		ddb.Updater(sop, _sop.Minus(sop.P)),
	)
	return err
}

var (
	_sop = ddb.UpdateFor[sop, []string]()
)

func encodeSOP(symbols Symbols, g curie.IRI, spock spock.SPOCK) sop {
	return sop{
		G:  "so|" + g,
		SO: encodeIV(symbols, spock.S, spock.O),
		P:  []string{symbols.ToSymbol(spock.P).String()},
	}
}

func decodeSOP(symbols Symbols, sop sop) []spock.SPOCK {
	seq := make([]spock.SPOCK, len(sop.P))
	s, o := decodeIV(symbols, sop.SO)

	for i, p := range sop.P {
		seq[i].S, seq[i].P, seq[i].O = s, symbols.FromString(p), o
	}

	return seq
}

//
// ⟨ Predicate, Object, Subject ⟩
//

type pos struct {
	G  curie.IRI `dynamodbav:"prefix"`
	PO string    `dynamodbav:"suffix"`
	S  []string  `dynamodbav:"s,stringset"`
}

func (pos pos) HashKey() curie.IRI { return pos.G }
func (pos pos) SortKey() curie.IRI { return curie.IRI(pos.PO) }

func (pos pos) ToSPOCK(symbols Symbols) []spock.SPOCK {
	return decodePOS(symbols, pos)
}

func (pos pos) Put(ctx context.Context, store *Store) error {
	_, err := store.pos.UpdateWith(ctx,
		ddb.Updater(pos, _pos.Union(pos.S)),
	)
	return err
}

func (pos pos) Cut(ctx context.Context, store *Store) error {
	_, err := store.pos.UpdateWith(ctx,
		ddb.Updater(pos, _pos.Minus(pos.S)),
	)
	return err
}

var (
	_pos = ddb.UpdateFor[pos, []string]()
)

func encodePOS(symbols Symbols, g curie.IRI, spock spock.SPOCK) pos {
	return pos{
		G:  "po|" + g,
		PO: encodeIV(symbols, spock.P, spock.O),
		S:  []string{symbols.ToSymbol(spock.S).String()},
	}
}

func decodePOS(symbols Symbols, pos pos) []spock.SPOCK {
	seq := make([]spock.SPOCK, len(pos.S))
	p, o := decodeIV(symbols, pos.PO)

	for i, s := range pos.S {
		seq[i].S, seq[i].P, seq[i].O = symbols.FromString(s), p, o
	}

	return seq
}

//
// ⟨ Predicate, Subject, Object ⟩
//

type pso struct {
	G  curie.IRI `dynamodbav:"prefix"`
	PS string    `dynamodbav:"suffix"`
	O  []string  `dynamodbav:"o,stringset"`
}

func (pso pso) HashKey() curie.IRI { return pso.G }
func (pso pso) SortKey() curie.IRI { return curie.IRI(pso.PS) }

func (pso pso) ToSPOCK(symbols Symbols) []spock.SPOCK {
	return decodePSO(symbols, pso)
}

func (pso pso) Put(ctx context.Context, store *Store) error {
	_, err := store.pso.UpdateWith(ctx,
		ddb.Updater(pso, _pso.Union(pso.O)),
	)
	return err
}

func (pso pso) Cut(ctx context.Context, store *Store) error {
	_, err := store.pso.UpdateWith(ctx,
		ddb.Updater(pso, _pso.Minus(pso.O)),
	)
	return err
}

var (
	_pso = ddb.UpdateFor[pso, []string]()
)

func encodePSO(symbols Symbols, g curie.IRI, spock spock.SPOCK) pso {
	return pso{
		G:  "ps|" + g,
		PS: encodeII(symbols, spock.P, spock.S),
		O:  []string{encodeValue(symbols, spock.O)},
	}
}

func decodePSO(symbols Symbols, pso pso) []spock.SPOCK {
	seq := make([]spock.SPOCK, len(pso.O))
	p, s := decodeII(symbols, pso.PS)

	for i, o := range pso.O {
		seq[i].S, seq[i].P, seq[i].O = s, p, decodeValue(symbols, o)
	}

	return seq
}

//
// ⟨ Object, Subject, Predicate ⟩
//

type osp struct {
	G  curie.IRI `dynamodbav:"prefix"`
	OS string    `dynamodbav:"suffix"`
	P  []string  `dynamodbav:"p,stringset"`
}

func (osp osp) HashKey() curie.IRI { return osp.G }
func (osp osp) SortKey() curie.IRI { return curie.IRI(osp.OS) }

func (osp osp) ToSPOCK(symbols Symbols) []spock.SPOCK {
	return decodeOSP(symbols, osp)
}

func (osp osp) Put(ctx context.Context, store *Store) error {
	_, err := store.osp.UpdateWith(ctx,
		ddb.Updater(osp, _osp.Union(osp.P)),
	)
	return err
}

func (osp osp) Cut(ctx context.Context, store *Store) error {
	_, err := store.osp.UpdateWith(ctx,
		ddb.Updater(osp, _osp.Minus(osp.P)),
	)
	return err
}

var (
	_osp = ddb.UpdateFor[osp, []string]()
)

func encodeOSP(symbols Symbols, g curie.IRI, spock spock.SPOCK) osp {
	return osp{
		G:  "os|" + g,
		OS: encodeVI(symbols, spock.O, spock.S),
		P:  []string{symbols.ToSymbol(spock.P).String()},
	}
}

func decodeOSP(symbols Symbols, osp osp) []spock.SPOCK {
	seq := make([]spock.SPOCK, len(osp.P))
	o, s := decodeVI(symbols, osp.OS)

	for i, p := range osp.P {
		seq[i].S, seq[i].P, seq[i].O = s, symbols.FromString(p), o
	}

	return seq
}

//
// ⟨ Object, Predicate, Subject ⟩
//

type ops struct {
	G  curie.IRI `dynamodbav:"prefix"`
	OP string    `dynamodbav:"suffix"`
	S  []string  `dynamodbav:"s,stringset"`
}

func (ops ops) HashKey() curie.IRI { return ops.G }
func (ops ops) SortKey() curie.IRI { return curie.IRI(ops.OP) }

func (ops ops) ToSPOCK(symbols Symbols) []spock.SPOCK {
	return decodeOPS(symbols, ops)
}

func (ops ops) Put(ctx context.Context, store *Store) error {
	_, err := store.ops.UpdateWith(ctx,
		ddb.Updater(ops, _ops.Union(ops.S)),
	)
	return err
}

func (ops ops) Cut(ctx context.Context, store *Store) error {
	_, err := store.ops.UpdateWith(ctx,
		ddb.Updater(ops, _ops.Minus(ops.S)),
	)
	return err
}

var (
	_ops = ddb.UpdateFor[ops, []string]()
)

func encodeOPS(symbols Symbols, g curie.IRI, spock spock.SPOCK) ops {
	return ops{
		G:  "op|" + g,
		OP: encodeVI(symbols, spock.O, spock.P),
		S:  []string{symbols.ToSymbol(spock.S).String()},
	}
}

func decodeOPS(symbols Symbols, ops ops) []spock.SPOCK {
	seq := make([]spock.SPOCK, len(ops.S))
	o, p := decodeVI(symbols, ops.OP)

	for i, s := range ops.S {
		seq[i].S, seq[i].P, seq[i].O = symbols.FromString(s), p, o
	}

	return seq
}
