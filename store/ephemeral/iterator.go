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
	"github.com/kshard/spock"
)

// evaluates query patterns against lists
type seqBuilder[A, B, C any] interface {
	L1(*skiplist.SkipList[A, *skiplist.SkipList[B, *skiplist.SkipList[C, k]]]) Seq[A, *skiplist.SkipList[B, *skiplist.SkipList[C, k]]]
	L2(*skiplist.SkipList[B, *skiplist.SkipList[C, k]]) Seq[B, *skiplist.SkipList[C, k]]
	L3(*skiplist.SkipList[C, k]) Seq[C, k]
	ToSPOCK(A, B, C) spock.SPOCK
}

type iterator[A, B, C any] struct {
	a   A
	b   B
	c   C
	abc Seq[A, *skiplist.SkipList[B, *skiplist.SkipList[C, k]]]
	_bc Seq[B, *skiplist.SkipList[C, k]]
	__c Seq[C, k]
	hlp seqBuilder[A, B, C]
}

func newIterator[A, B, C any](
	hlp seqBuilder[A, B, C],
	seq *skiplist.SkipList[A, *skiplist.SkipList[B, *skiplist.SkipList[C, k]]],
) *iterator[A, B, C] {
	return &iterator[A, B, C]{
		hlp: hlp,
		abc: hlp.L1(seq),
	}
}

func (iter *iterator[A, B, C]) Head() spock.SPOCK {
	return iter.hlp.ToSPOCK(iter.a, iter.b, iter.c)
}

func (iter *iterator[A, B, C]) Next() bool {
	if iter._bc == nil {
		if iter.abc == nil || !iter.abc.Next() {
			return false
		}
		a, _bc := iter.abc.Head()
		iter.a = a
		iter._bc = iter.hlp.L2(_bc)
	}

	if iter.__c == nil {
		if iter._bc == nil || !iter._bc.Next() {
			iter._bc = nil
			return iter.Next()
		}

		b, __c := iter._bc.Head()
		iter.b = b
		iter.__c = iter.hlp.L3(__c)
	}

	if iter.__c == nil || !iter.__c.Next() {
		iter.__c = nil
		return iter.Next()
	}

	iter.c, _ = iter.__c.Head()

	return true
}

func (iter *iterator[A, B, C]) FMap(f func(spock.SPOCK) error) error {
	for iter.Next() {
		if err := f(iter.Head()); err != nil {
			return err
		}
	}
	return nil
}
