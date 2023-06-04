package hybrid

import (
	"github.com/fogfish/skiplist"
	"github.com/kshard/spock"
)

type seqBuilder[A, B, C any] interface {
	L0(*skiplist.SkipList[h, *skiplist.SkipList[A, *skiplist.SkipList[B, *skiplist.SkipList[C, k]]]]) Seq[h, *skiplist.SkipList[A, *skiplist.SkipList[B, *skiplist.SkipList[C, k]]]]
	L1(*skiplist.SkipList[A, *skiplist.SkipList[B, *skiplist.SkipList[C, k]]]) Seq[A, *skiplist.SkipList[B, *skiplist.SkipList[C, k]]]
	L2(*skiplist.SkipList[B, *skiplist.SkipList[C, k]]) Seq[B, *skiplist.SkipList[C, k]]
	L3(*skiplist.SkipList[C, k]) Seq[C, k]
	Swap(h) *skiplist.SkipList[A, *skiplist.SkipList[B, *skiplist.SkipList[C, k]]]
	ToSPOCK(A, B, C) spock.SPOCK
}

type iterator[A, B, C any] struct {
	a       A
	b       B
	c       C
	seq     Seq[h, *skiplist.SkipList[A, *skiplist.SkipList[B, *skiplist.SkipList[C, k]]]]
	abc     Seq[A, *skiplist.SkipList[B, *skiplist.SkipList[C, k]]]
	_bc     Seq[B, *skiplist.SkipList[C, k]]
	__c     Seq[C, k]
	builder seqBuilder[A, B, C]
}

func newIterator[A, B, C any](
	builder seqBuilder[A, B, C],
	seq *skiplist.SkipList[h, *skiplist.SkipList[A, *skiplist.SkipList[B, *skiplist.SkipList[C, k]]]],
) *iterator[A, B, C] {
	return &iterator[A, B, C]{
		seq:     builder.L0(seq),
		builder: builder,
	}
}

func (iter *iterator[A, B, C]) Head() spock.SPOCK {
	return iter.builder.ToSPOCK(iter.a, iter.b, iter.c)
}

func (iter *iterator[A, B, C]) Next() bool {
	if iter.abc == nil {
		if iter.seq == nil || !iter.seq.Next() {
			return false
		}

		h, abc := iter.seq.Head()
		if abc == nil {
			abc = iter.builder.Swap(h)
		}
		iter.abc = iter.builder.L1(abc)

	}

	if iter._bc == nil {
		if iter.abc == nil || !iter.abc.Next() {
			return false
		}
		a, _bc := iter.abc.Head()
		iter.a = a
		iter._bc = iter.builder.L2(_bc)
	}

	if iter.__c == nil {
		if iter._bc == nil || !iter._bc.Next() {
			iter._bc = nil
			return iter.Next()
		}

		b, __c := iter._bc.Head()
		iter.b = b
		iter.__c = iter.builder.L3(__c)
	}

	if iter.__c == nil || !iter.__c.Next() {
		iter.__c = nil
		return iter.Next()
	}

	iter.c, _ = iter.__c.Head()

	return true
}

func (iter *iterator[A, B, C]) FMap(f func(spock.SPOCK) error) error {
	return nil
}
