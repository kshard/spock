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
	"fmt"

	"github.com/fogfish/curie"
	"github.com/kshard/spock"
)

type notSupported struct{ spock.Pattern }

func (err notSupported) Error() string { return fmt.Sprintf("not supported %s", err.Pattern.Dump()) }
func (notSupported) NotSupported()     {}

func (store *Store) streamSPO(ctx context.Context, graph curie.IRI, q spock.Pattern) (spock.Stream, error) {
	key := spo{G: "sp|" + graph}

	switch {
	case q.HintForS == spock.HINT_MATCH && q.HintForP == spock.HINT_NONE:
		key.SP = encodeI(store.iri, q.S.Value)
	case q.HintForS == spock.HINT_MATCH && q.HintForP == spock.HINT_MATCH:
		key.SP = encodeII(store.iri, q.S.Value, q.P.Value)
	case q.HintForS == spock.HINT_MATCH && q.HintForP == spock.HINT_FILTER_PREFIX:
		key.SP = encodeII(store.iri, q.S.Value, q.P.Value)
	case q.HintForS == spock.HINT_FILTER_PREFIX && q.HintForP == spock.HINT_NONE:
		key.SP = encodeI(store.iri, q.S.Value)
	default:
		return nil, &notSupported{q}
	}

	var stream spock.Stream = &Unfold[spo]{
		symbols: store.iri,
		seq:     NewIterator(store.spo, key),
	}

	if q.O != nil {
		stream = spock.NewFilterO(q.HintForO, q.O, stream)
	}

	return stream, nil
}

func (store *Store) streamSOP(ctx context.Context, graph curie.IRI, q spock.Pattern) (spock.Stream, error) {
	key := sop{G: "so|" + graph}

	switch {
	case q.HintForS == spock.HINT_MATCH && q.HintForO == spock.HINT_NONE:
		key.SO = encodeI(store.iri, q.S.Value)
	case q.HintForS == spock.HINT_MATCH && q.HintForO == spock.HINT_MATCH:
		key.SO = encodeIV(store.iri, q.S.Value, q.O.Value)
	case q.HintForS == spock.HINT_MATCH && q.HintForO == spock.HINT_FILTER_PREFIX:
		key.SO = encodeIV(store.iri, q.S.Value, q.O.Value)
	case q.HintForS == spock.HINT_FILTER_PREFIX && q.HintForO == spock.HINT_NONE:
		key.SO = encodeI(store.iri, q.S.Value)
	default:
		return nil, &notSupported{q}
	}

	var stream spock.Stream = &Unfold[sop]{
		symbols: store.iri,
		seq:     NewIterator(store.sop, key),
	}

	if q.P != nil {
		stream = spock.NewFilterP(q.HintForP, q.P, stream)
	}

	return stream, nil
}

func (store *Store) streamPSO(ctx context.Context, graph curie.IRI, q spock.Pattern) (spock.Stream, error) {
	key := pso{G: "ps|" + graph}

	switch {
	case q.HintForP == spock.HINT_MATCH && q.HintForS == spock.HINT_NONE:
		key.PS = encodeI(store.iri, q.P.Value)
	case q.HintForP == spock.HINT_MATCH && q.HintForS == spock.HINT_MATCH:
		key.PS = encodeII(store.iri, q.P.Value, q.S.Value)
	case q.HintForP == spock.HINT_MATCH && q.HintForS == spock.HINT_FILTER_PREFIX:
		key.PS = encodeII(store.iri, q.P.Value, q.S.Value)
	case q.HintForP == spock.HINT_FILTER_PREFIX && q.HintForS == spock.HINT_NONE:
		key.PS = encodeI(store.iri, q.P.Value)
	default:
		return nil, &notSupported{q}
	}

	var stream spock.Stream = &Unfold[pso]{
		symbols: store.iri,
		seq:     NewIterator(store.pso, key),
	}

	if q.O != nil {
		stream = spock.NewFilterO(q.HintForO, q.O, stream)
	}

	return stream, nil
}

func (store *Store) streamPOS(ctx context.Context, graph curie.IRI, q spock.Pattern) (spock.Stream, error) {
	key := pos{G: "po|" + graph}

	switch {
	case q.HintForP == spock.HINT_MATCH && q.HintForO == spock.HINT_NONE:
		key.PO = encodeI(store.iri, q.P.Value)
	case q.HintForP == spock.HINT_MATCH && q.HintForO == spock.HINT_MATCH:
		key.PO = encodeIV(store.iri, q.P.Value, q.O.Value)
	case q.HintForP == spock.HINT_MATCH && q.HintForO == spock.HINT_FILTER_PREFIX:
		key.PO = encodeIV(store.iri, q.P.Value, q.O.Value)
	case q.HintForP == spock.HINT_FILTER_PREFIX && q.HintForO == spock.HINT_NONE:
		key.PO = encodeI(store.iri, q.P.Value)
	default:
		return nil, &notSupported{q}
	}

	var stream spock.Stream = &Unfold[pos]{
		symbols: store.iri,
		seq:     NewIterator(store.pos, key),
	}

	if q.S != nil {
		stream = spock.NewFilterS(q.HintForS, q.S, stream)
	}

	return stream, nil
}

func (store *Store) streamOSP(ctx context.Context, graph curie.IRI, q spock.Pattern) (spock.Stream, error) {
	key := osp{G: "os|" + graph}

	switch {
	case q.HintForO == spock.HINT_MATCH && q.HintForS == spock.HINT_NONE:
		key.OS = encodeValue(store.iri, q.O.Value)
	case q.HintForO == spock.HINT_MATCH && q.HintForS == spock.HINT_MATCH:
		key.OS = encodeVI(store.iri, q.O.Value, q.S.Value)
	case q.HintForO == spock.HINT_MATCH && q.HintForS == spock.HINT_FILTER_PREFIX:
		key.OS = encodeVI(store.iri, q.O.Value, q.S.Value)
	case q.HintForO == spock.HINT_FILTER_PREFIX && q.HintForS == spock.HINT_NONE:
		key.OS = encodeValue(store.iri, q.O.Value)
	default:
		return nil, &notSupported{q}
	}

	var stream spock.Stream = &Unfold[osp]{
		symbols: store.iri,
		seq:     NewIterator(store.osp, key),
	}

	if q.P != nil {
		stream = spock.NewFilterP(q.HintForP, q.P, stream)
	}

	return stream, nil
}

func (store *Store) streamOPS(ctx context.Context, graph curie.IRI, q spock.Pattern) (spock.Stream, error) {
	key := ops{G: "op|" + graph}

	switch {
	case q.HintForO == spock.HINT_MATCH && q.HintForP == spock.HINT_NONE:
		key.OP = encodeValue(store.iri, q.O.Value)
	case q.HintForO == spock.HINT_MATCH && q.HintForP == spock.HINT_MATCH:
		key.OP = encodeVI(store.iri, q.O.Value, q.P.Value)
	case q.HintForO == spock.HINT_MATCH && q.HintForP == spock.HINT_FILTER_PREFIX:
		key.OP = encodeVI(store.iri, q.O.Value, q.P.Value)
	case q.HintForO == spock.HINT_FILTER_PREFIX && q.HintForP == spock.HINT_NONE:
		key.OP = encodeValue(store.iri, q.O.Value)
	default:
		return nil, &notSupported{q}
	}

	var stream spock.Stream = &Unfold[ops]{
		symbols: store.iri,
		seq:     NewIterator(store.ops, key),
	}

	if q.S != nil {
		stream = spock.NewFilterS(q.HintForS, q.S, stream)
	}

	return stream, nil
}
