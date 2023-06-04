package hybrid

import (
	"fmt"

	"github.com/fogfish/faults"
	"github.com/fogfish/golem/trait/seq"
	"github.com/fogfish/segment"
	"github.com/fogfish/skiplist"
	"github.com/kshard/spock"
	"github.com/kshard/xsd"
)

type Store struct {
	size int
	spo  *segment.Map[sp, __o]
	sop  *segment.Map[so, __p]
	pso  *segment.Map[ps, __o]
	pos  *segment.Map[po, __s]
	osp  *segment.Map[os, __p]
	ops  *segment.Map[op, __s]
}

// Create new instance of knowledge storage
func New(capacity int) (*Store, error) {
	dspo := NewDynamo[xsd.Symbol]("spo")
	spo, err := segment.New[sp, __o](capacity, nil, dspo, dspo)
	if err != nil {
		return nil, err
	}

	dsop := NewDynamo[xsd.AnyURI]("sop")
	sop, err := segment.New[so, __p](capacity, nil, dsop, dsop)
	if err != nil {
		return nil, err
	}

	dpso := NewDynamo[xsd.Symbol]("pso")
	pso, err := segment.New[ps, __o](capacity, nil, dpso, dpso)
	if err != nil {
		return nil, err
	}

	dpos := NewDynamo[xsd.AnyURI]("pos")
	pos, err := segment.New[po, __s](capacity, nil, dpos, dpos)
	if err != nil {
		return nil, err
	}

	dosp := NewDynamo[xsd.AnyURI]("osp")
	osp, err := segment.New[os, __p](capacity, nil, dosp, dosp)
	if err != nil {
		return nil, err
	}

	dops := NewDynamo[xsd.AnyURI]("ops")
	ops, err := segment.New[op, __s](capacity, nil, dops, dops)
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

func (store *Store) Sync() error {
	store.spo.Sync()
	store.sop.Sync()
	store.pso.Sync()
	store.pos.Sync()
	store.osp.Sync()
	store.ops.Sync()
	return nil
}

// Size returns number of knowledge statements in the store
func (store *Store) Size() int {
	return store.size
}

func (store *Store) Add(bag spock.Bag) {
	for _, spock := range bag {
		store.Put(spock)
	}
}

func (store *Store) Put(spock spock.SPOCK) {
	// spock.K = guid.L(guid.Clock)

	store.putSuffixO(spock)
	store.putSuffixS(spock)
	store.putSuffixP(spock)

	store.size++
}

func (store *Store) putSuffixO(spock spock.SPOCK) bool {
	sp := uint64(spock.S)<<32 | uint64(spock.P)
	ps := uint64(spock.P)<<32 | uint64(spock.S)
	__o, err := store.spo.Get(sp)
	if faults.IsNotFound(err) {
		__o = skiplist.NewSet[xsd.Symbol]()
		// TODO: error
		store.spo.Put(sp, __o)
		store.pso.Put(ps, __o)
	}

	return __o.Add(spock.O)
}

func (store *Store) putSuffixS(spock spock.SPOCK) bool {
	po := uint64(spock.P)<<32 | uint64(spock.O)
	op := uint64(spock.O)<<32 | uint64(spock.P)
	__s, err := store.pos.Get(po)
	if faults.IsNotFound(err) {
		__s = skiplist.NewSet[xsd.AnyURI]()
		// TODO: error
		store.pos.Put(po, __s)
		store.ops.Put(op, __s)
	}

	return __s.Add(spock.S)
}

func (store *Store) putSuffixP(spock spock.SPOCK) bool {
	so := uint64(spock.S)<<32 | uint64(spock.O)
	os := uint64(spock.O)<<32 | uint64(spock.S)
	__p, err := store.sop.Get(so)
	if faults.IsNotFound(err) {
		__p = skiplist.NewSet[xsd.AnyURI]()
		// TODO: error
		store.sop.Put(so, __p)
		store.osp.Put(os, __p)
	}

	return __p.Add(spock.P)
}

func (store *Store) Match(q spock.Pattern) (seq.Seq[spock.SPOCK], error) {
	if q.HintForS != spock.HINT_MATCH && q.HintForS != spock.HINT_NONE {
		return nil, &notSupported{q}
	}

	if q.HintForP != spock.HINT_MATCH && q.HintForP != spock.HINT_NONE {
		return nil, &notSupported{q}
	}

	if q.HintForO != spock.HINT_MATCH && q.HintForO != spock.HINT_NONE { //&& q.O.Value.XSDType() == xsd.XSD_ANYURI {
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
