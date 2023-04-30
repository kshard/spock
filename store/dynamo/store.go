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
	"github.com/fogfish/dynamo/v2"
	"github.com/fogfish/dynamo/v2/service/ddb"
	"github.com/kshard/spock"
)

type Store struct {
	spo *ddb.Storage[spo]
	sop *ddb.Storage[sop]
	pso *ddb.Storage[pso]
	pos *ddb.Storage[pos]
	osp *ddb.Storage[osp]
	ops *ddb.Storage[ops]
}

func New(connector string, opts ...dynamo.Option) (*Store, error) {
	spo, err := ddb.New[spo](connector, opts...)
	if err != nil {
		return nil, err
	}

	sop, err := ddb.New[sop](connector, opts...)
	if err != nil {
		return nil, err
	}

	pso, err := ddb.New[pso](connector, opts...)
	if err != nil {
		return nil, err
	}

	pos, err := ddb.New[pos](connector, opts...)
	if err != nil {
		return nil, err
	}

	osp, err := ddb.New[osp](connector, opts...)
	if err != nil {
		return nil, err
	}

	ops, err := ddb.New[ops](connector, opts...)
	if err != nil {
		return nil, err
	}

	return &Store{
		spo: spo,
		sop: sop,
		pso: pso,
		pos: pos,
		osp: osp,
		ops: ops,
	}, nil
}

func Add(ctx context.Context, store *Store, graph curie.IRI, bag spock.Bag) (spock.Bag, error) {
	for i, spock := range bag {
		if err := Put(ctx, store, graph, spock); err != nil {
			return bag[i:], err
		}
	}

	return nil, nil
}

func Put(ctx context.Context, store *Store, graph curie.IRI, spock spock.SPOCK) error {
	seq := []Writer{
		encodeSPO(graph, spock),
		encodeSOP(graph, spock),
		encodePOS(graph, spock),
		encodePSO(graph, spock),
		encodeOPS(graph, spock),
		encodeOSP(graph, spock),
	}

	for i := 0; i < len(seq); i++ {
		if err := seq[i].Put(ctx, store); err != nil {
			for k := 0; k < i; k++ {
				if err := seq[k].Cut(ctx, store); err != nil {
					// TODO: log error
				}
			}
			return err
		}
	}

	return nil
}

func Match(ctx context.Context, store *Store, graph curie.IRI, q spock.Pattern) (spock.Stream, error) {
	switch q.Strategy {
	case spock.STRATEGY_SPO:
		return store.streamSPO(ctx, graph, q)
	case spock.STRATEGY_SOP:
		return store.streamSOP(ctx, graph, q)
	case spock.STRATEGY_PSO:
		return store.streamPSO(ctx, graph, q)
	case spock.STRATEGY_POS:
		return store.streamPOS(ctx, graph, q)
	case spock.STRATEGY_OSP:
		return store.streamOSP(ctx, graph, q)
	case spock.STRATEGY_OPS:
		return store.streamOPS(ctx, graph, q)
	default:
		panic(fmt.Errorf("unknown strategy"))
	}
}
