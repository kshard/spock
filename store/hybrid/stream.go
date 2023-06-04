package hybrid

import (
	"fmt"

	"github.com/fogfish/faults"
	"github.com/fogfish/golem/trait/pair"
	"github.com/fogfish/golem/trait/seq"
	"github.com/fogfish/segment"
	"github.com/fogfish/skiplist"
	"github.com/kshard/spock"
	"github.com/kshard/xsd"
)

type notSupported struct{ spock.Pattern[xsd.AnyURI] }

func (err notSupported) Error() string { return fmt.Sprintf("not supported %s", err.Pattern.Dump()) }
func (notSupported) NotSupported()     {}

type X[C skiplist.Key] interface {
	L1() (uint64, uint64)
	L2() uint64
	ToSPOCK(uint64, seq.Seq[C]) seq.Seq[spock.SPOCK]
}

// helper function to query prefix of the index
func queryIRI[A, B, C skiplist.Key](
	qA *spock.Predicate[A],
	qB *spock.Predicate[B],
	qC *spock.Predicate[C],
	lv X[C],
	list *segment.Map[uint64, *skiplist.Set[C]],
) seq.Seq[spock.SPOCK] {
	if qA != nil && qA.Clause == spock.EQ && qB != nil && qB.Clause == spock.EQ {
		ab := lv.L2()
		__x, err := list.Get(ab)
		if faults.IsNotFound(err) {
			return nil
		}

		return lv.ToSPOCK(ab, queryXSD(qC, __x))
	}

	if qA != nil && qA.Clause == spock.EQ {
		ab, ab1 := lv.L1()
		a__, err := list.Successor(ab)
		if err != nil {
			return nil
		}
		a__ = pair.TakeWhile(a__,
			func(ab uint64, __x *skiplist.Set[C]) bool { return ab < ab1 },
		)
		return pair.ToSeq(a__,
			func(ab uint64, __x *skiplist.Set[C]) seq.Seq[spock.SPOCK] {
				return lv.ToSPOCK(ab, queryXSD(qC, __x))
			},
		)
	}

	a__, err := list.Values()
	if err != nil {
		return nil
	}
	return pair.ToSeq(a__,
		func(ab uint64, __x *skiplist.Set[C]) seq.Seq[spock.SPOCK] {
			return lv.ToSPOCK(ab, queryXSD(qC, __x))
		},
	)
}

func queryXSD[A skiplist.Key](
	pred *spock.Predicate[A],
	list *skiplist.Set[A],
) seq.Seq[A] {
	if pred != nil && pred.Clause == spock.EQ {
		if has, _ := list.Has(pred.Value); has {
			return seq.From(pred.Value)
		}
		return nil
	}

	return skiplist.ForSet(list, list.Values())
}

func (store *Store) streamSPO(q spock.Pattern[xsd.AnyURI]) (seq.Seq[spock.SPOCK], error) {
	return queryIRI[xsd.AnyURI, xsd.AnyURI, xsd.AnyURI](
		q.S, q.P, q.O, querySPO(q), store.spo,
	), nil
}

func (store *Store) streamSOP(q spock.Pattern[xsd.AnyURI]) (seq.Seq[spock.SPOCK], error) {
	return queryIRI[xsd.AnyURI, xsd.AnyURI, xsd.AnyURI](
		q.S, q.O, q.P, querySOP(q), store.sop,
	), nil
}

func (store *Store) streamPSO(q spock.Pattern[xsd.AnyURI]) (seq.Seq[spock.SPOCK], error) {
	return queryIRI[xsd.AnyURI, xsd.AnyURI, xsd.AnyURI](
		q.P, q.S, q.O, queryPSO(q), store.pso,
	), nil
}

func (store *Store) streamPOS(q spock.Pattern[xsd.AnyURI]) (seq.Seq[spock.SPOCK], error) {
	return queryIRI[xsd.AnyURI, xsd.AnyURI, xsd.AnyURI](
		q.P, q.O, q.S, queryPOS(q), store.pos,
	), nil
}

func (store *Store) streamOSP(q spock.Pattern[xsd.AnyURI]) (seq.Seq[spock.SPOCK], error) {
	return queryIRI[xsd.AnyURI, xsd.AnyURI, xsd.AnyURI](
		q.O, q.S, q.P, queryOSP(q), store.osp,
	), nil
}

func (store *Store) streamOPS(q spock.Pattern[xsd.AnyURI]) (seq.Seq[spock.SPOCK], error) {
	return queryIRI[xsd.AnyURI, xsd.AnyURI, xsd.AnyURI](
		q.O, q.P, q.S, queryOPS(q), store.ops,
	), nil
}

type querySPO spock.Pattern[xsd.AnyURI]

func (q querySPO) L1() (uint64, uint64) {
	return uint64(q.S.Value) << 32, uint64(q.S.Value+1) << 32
}

func (q querySPO) L2() uint64 {
	return uint64(q.S.Value)<<32 | uint64(q.P.Value)
}

func (q querySPO) ToSPOCK(sp uint64, o seq.Seq[xsd.AnyURI]) seq.Seq[spock.SPOCK] {
	return seq.Map(o,
		func(o xsd.AnyURI) spock.SPOCK {
			return spock.SPOCK{
				S: xsd.AnyURI(sp >> 32),
				P: xsd.AnyURI(sp & (0xffffffff)),
				O: o,
			}
		},
	)
}

type querySOP spock.Pattern[xsd.AnyURI]

func (q querySOP) L1() (uint64, uint64) {
	return uint64(q.S.Value) << 32, uint64(q.S.Value+1) << 32
}

func (q querySOP) L2() uint64 {
	return uint64(q.S.Value)<<32 | uint64(q.O.Value)
}

func (q querySOP) ToSPOCK(so uint64, p seq.Seq[xsd.AnyURI]) seq.Seq[spock.SPOCK] {
	return seq.Map(p,
		func(p xsd.AnyURI) spock.SPOCK {
			return spock.SPOCK{
				S: xsd.AnyURI(so >> 32),
				P: p,
				O: xsd.AnyURI(so & (0xffffffff)),
			}
		},
	)
}

type queryPSO spock.Pattern[xsd.AnyURI]

func (q queryPSO) L1() (uint64, uint64) {
	return uint64(q.P.Value) << 32, uint64(q.P.Value+1) << 32
}

func (q queryPSO) L2() uint64 {
	return uint64(q.P.Value)<<32 | uint64(q.S.Value)
}

func (q queryPSO) ToSPOCK(ps uint64, o seq.Seq[xsd.AnyURI]) seq.Seq[spock.SPOCK] {
	return seq.Map(o,
		func(o xsd.AnyURI) spock.SPOCK {
			return spock.SPOCK{
				S: xsd.AnyURI(ps & (0xffffffff)),
				P: xsd.AnyURI(ps >> 32),
				O: o,
			}
		},
	)
}

type queryPOS spock.Pattern[xsd.AnyURI]

func (q queryPOS) L1() (uint64, uint64) {
	return uint64(q.P.Value) << 32, uint64(q.P.Value+1) << 32
}

func (q queryPOS) L2() uint64 {
	return uint64(q.P.Value)<<32 | uint64(q.O.Value)
}

func (q queryPOS) ToSPOCK(po uint64, s seq.Seq[xsd.AnyURI]) seq.Seq[spock.SPOCK] {
	return seq.Map(s,
		func(s xsd.AnyURI) spock.SPOCK {
			return spock.SPOCK{
				S: s,
				P: xsd.AnyURI(po >> 32),
				O: xsd.AnyURI(po & (0xffffffff)),
			}
		},
	)
}

type queryOSP spock.Pattern[xsd.AnyURI]

func (q queryOSP) L1() (uint64, uint64) {
	return uint64(q.O.Value) << 32, uint64(q.O.Value+1) << 32
}

func (q queryOSP) L2() uint64 {
	return uint64(q.O.Value)<<32 | uint64(q.S.Value)
}

func (q queryOSP) ToSPOCK(os uint64, p seq.Seq[xsd.AnyURI]) seq.Seq[spock.SPOCK] {
	return seq.Map(p,
		func(p xsd.AnyURI) spock.SPOCK {
			return spock.SPOCK{
				S: xsd.AnyURI(os & (0xffffffff)),
				P: p,
				O: xsd.AnyURI(os >> 32),
			}
		},
	)
}

type queryOPS spock.Pattern[xsd.AnyURI]

func (q queryOPS) L1() (uint64, uint64) {
	return uint64(q.O.Value) << 32, uint64(q.O.Value+1) << 32
}

func (q queryOPS) L2() uint64 {
	return uint64(q.O.Value)<<32 | uint64(q.P.Value)
}

func (q queryOPS) ToSPOCK(op uint64, s seq.Seq[xsd.AnyURI]) seq.Seq[spock.SPOCK] {
	return seq.Map(s,
		func(s xsd.AnyURI) spock.SPOCK {
			return spock.SPOCK{
				S: s,
				P: xsd.AnyURI(op & (0xffffffff)),
				O: xsd.AnyURI(op >> 32),
			}
		},
	)
}
