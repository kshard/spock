package hybrid

import (
	"github.com/fogfish/skiplist"
	"github.com/kshard/xsd"
)

type (
	s  = uint32
	p  = uint32
	o  = uint32
	sp = uint64
	so = uint64
	ps = uint64
	po = uint64
	os = uint64
	op = uint64
)

// index types for 3rd faction
type __s = *skiplist.Set[xsd.AnyURI]
type __p = *skiplist.Set[xsd.AnyURI]
type __o = *skiplist.Set[xsd.AnyURI]
