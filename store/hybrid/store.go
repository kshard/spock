package hybrid

import (
	"fmt"

	"github.com/fogfish/faults"
	"github.com/fogfish/golem/trait/seq"
	"github.com/fogfish/segment"
	"github.com/fogfish/segment/datastore/dynamodb"
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
func New() (*Store, error) {
	dspo, _ := dynamodb.New[sp, __o]("spo", "ddb:///segment")
	spo, err := segment.New[sp, __o](dspo, dspo,
		skiplist.MapWithBlockSize[sp, __o](1024),
	)
	if err != nil {
		return nil, err
	}

	dsop, _ := dynamodb.New[so, __p]("sop", "ddb:///segment")
	sop, err := segment.New[so, __p](dsop, dsop,
		skiplist.MapWithBlockSize[so, __p](1024),
	)
	if err != nil {
		return nil, err
	}

	dpso, _ := dynamodb.New[ps, __o]("pso", "ddb:///segment")
	pso, err := segment.New[ps, __o](dpso, dpso,
		skiplist.MapWithBlockSize[ps, __o](1024),
	)
	if err != nil {
		return nil, err
	}

	dpos, _ := dynamodb.New[po, __s]("pos", "ddb:///segment")
	pos, err := segment.New[po, __s](dpos, dpos,
		skiplist.MapWithBlockSize[po, __s](1024),
	)
	if err != nil {
		return nil, err
	}

	dosp, _ := dynamodb.New[os, __p]("osp", "ddb:///segment")
	osp, err := segment.New[os, __p](dosp, dosp,
		skiplist.MapWithBlockSize[os, __p](1024),
	)
	if err != nil {
		return nil, err
	}

	dops, _ := dynamodb.New[op, __s]("ops", "ddb:///segment")
	ops, err := segment.New[op, __s](dops, dops,
		skiplist.MapWithBlockSize[op, __s](1024),
	)
	if err != nil {
		return nil, err
	}

	// TODO: only 1/2 of indexes must be persistent

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

func (store *Store) Put(spock spock.SPOCK) error {
	// spock.K = guid.L(guid.Clock)

	kind := spock.O.XSDType()
	if kind != xsd.XSD_ANYURI /*|| kind != xsd.XSD_SYMBOL*/ {
		return fmt.Errorf("not supported %T", spock.O)
	}

	o, ok := spock.O.(xsd.AnyURI)
	if !ok {
		return fmt.Errorf("not supported %T", spock.O)
	}

	store.putSuffixO(spock.S, spock.P, o)
	store.putSuffixS(spock.S, spock.P, o)
	store.putSuffixP(spock.S, spock.P, o)

	store.size++
	return nil
}

func (store *Store) putSuffixO(s, p, o xsd.AnyURI) bool {
	sp := uint64(s)<<32 | uint64(p)
	ps := uint64(p)<<32 | uint64(s)
	__o, err := store.spo.Get(sp)
	if faults.IsNotFound(err) {
		__o = skiplist.NewSet[xsd.AnyURI]()
		// TODO: error
		store.spo.Put(sp, __o)
		store.pso.Put(ps, __o)
	}

	has, _ := __o.Add(o)
	return has
}

func (store *Store) putSuffixS(s, p, o xsd.AnyURI) bool {
	po := uint64(p)<<32 | uint64(o)
	op := uint64(o)<<32 | uint64(p)
	__s, err := store.pos.Get(po)
	if faults.IsNotFound(err) {
		__s = skiplist.NewSet[xsd.AnyURI]()
		// TODO: error
		store.pos.Put(po, __s)
		store.ops.Put(op, __s)
	}

	has, _ := __s.Add(s)
	return has
}

func (store *Store) putSuffixP(s, p, o xsd.AnyURI) bool {
	so := uint64(s)<<32 | uint64(o)
	os := uint64(o)<<32 | uint64(s)
	__p, err := store.sop.Get(so)
	if faults.IsNotFound(err) {
		__p = skiplist.NewSet[xsd.AnyURI]()
		// TODO: error
		store.sop.Put(so, __p)
		store.osp.Put(os, __p)
	}

	has, _ := __p.Add(p)
	return has
}

func (store *Store) Match(q spock.Pattern[xsd.AnyURI]) (seq.Seq[spock.SPOCK], error) {
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
