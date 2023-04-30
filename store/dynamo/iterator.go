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

	"github.com/fogfish/dynamo/v2"
	"github.com/fogfish/dynamo/v2/service/ddb"
	"github.com/kshard/spock"
)

type none string

func (none) MatchOpt() {}

type Seq[T dynamo.Thing] interface {
	Head() T
	Next() bool
}

func NewIterator[T dynamo.Thing](store *ddb.Storage[T], query T) Seq[T] {
	return &Iterator[T]{
		store:  store,
		query:  query,
		cursor: none(""),
	}
}

type Iterator[T dynamo.Thing] struct {
	store  *ddb.Storage[T]
	query  T
	cursor dynamo.MatchOpt
	seq    []T
}

func (iter *Iterator[T]) Head() T {
	return iter.seq[0]
}

func (iter *Iterator[T]) Next() bool {
	if iter.seq != nil && len(iter.seq) > 1 {
		iter.seq = iter.seq[1:]
		return true
	}

	if iter.cursor == nil {
		return false
	}

	var err error
	iter.seq, iter.cursor, err = iter.store.Match(context.TODO(),
		iter.query, iter.cursor, dynamo.Limit(2),
	)
	if err != nil {
		return false
	}

	if len(iter.seq) == 0 {
		return false
	}

	return true
}

type Unfold[T dynamo.Thing] struct {
	seq Seq[T]
	bag []spock.SPOCK
}

func (unfold *Unfold[T]) Head() spock.SPOCK {
	return unfold.bag[0]
}

func (unfold *Unfold[T]) Next() bool {
	if unfold.bag != nil && len(unfold.bag) > 1 {
		unfold.bag = unfold.bag[1:]
		return true
	}

	if !unfold.seq.Next() {
		return false
	}

	switch vv := any(unfold.seq.Head()).(type) {
	case interface{ ToSPOCK() []spock.SPOCK }:
		unfold.bag = vv.ToSPOCK()
	}

	return true
}

func (unfold *Unfold[T]) FMap(f func(spock.SPOCK) error) error {
	for unfold.Next() {
		if err := f(unfold.Head()); err != nil {
			return err
		}
	}
	return nil
}
