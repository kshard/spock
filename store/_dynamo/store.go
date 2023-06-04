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
	"github.com/kshard/spock/internal/symbol"
	"github.com/kshard/xsd"
)

type Symbols interface {
	ToSymbol(xsd.AnyURI) symbol.Symbol
	ToAnyURI(symbol.Symbol) xsd.AnyURI
	FromString(string) xsd.AnyURI
	Write(context.Context) error
}

func NewSymbols(namespace curie.IRI, connector string, opts ...dynamo.Option) (Symbols, error) {
	store, err := symbol.NewStore(connector, opts...)
	if err != nil {
		return nil, err
	}

	symbols := symbol.New(namespace, store)

	// if err := symbols.Read(context.Background()); err != nil {
	// 	return nil, err
	// }

	return symbols, nil
}

type Store struct {
	iri Symbols
	spo *ddb.Storage[spo]
	sop *ddb.Storage[sop]
	pso *ddb.Storage[pso]
	pos *ddb.Storage[pos]
	osp *ddb.Storage[osp]
	ops *ddb.Storage[ops]
}

func New(iri Symbols, connector string, opts ...dynamo.Option) (*Store, error) {
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
		iri: iri,
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
		// t := time.Now()
		if err := Put(ctx, store, graph, spock); err != nil {
			return bag[i:], err
		}
		// fmt.Println(time.Since(t))
	}

	return nil, nil
}

func Put(ctx context.Context, store *Store, graph curie.IRI, spock spock.SPOCK) error {
	seq := []Writer{
		encodeSPO(store.iri, graph, spock),
		encodeSOP(store.iri, graph, spock),
		encodePOS(store.iri, graph, spock),
		encodePSO(store.iri, graph, spock),
		encodeOPS(store.iri, graph, spock),
		encodeOSP(store.iri, graph, spock),
	}

	sem := make(chan struct{}, len(seq))

	for i := 0; i < len(seq); i++ {
		// t := time.Now()
		go func(id int) {
			seq[id].Put(ctx, store)
			sem <- struct{}{}
			// if err := seq[i].Put(ctx, store); err != nil {
			// 	for k := 0; k < i; k++ {
			// 		if err := seq[k].Cut(ctx, store); err != nil {
			// 			// TODO: log error
			// 		}
			// 	}
			// 	return err
			// }
		}(i)
		// fmt.Println(time.Since(t))
	}

	// t := time.Now()
	for i := 0; i < len(seq); i++ {
		<-sem
	}
	// fmt.Println(time.Since(t))

	return nil
}

func Match(ctx context.Context, store *Store, graph curie.IRI, q spock.Pattern) (spock.Stream, error) {
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
