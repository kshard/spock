package hybrid

import (
	"github.com/fogfish/skiplist"
	"github.com/kshard/xsd"
)

// components of <s,p,o,c,k> triple
type h = uint32
type s = xsd.AnyURI // subject
type p = xsd.AnyURI // predicate
type o = xsd.Value  // object
type c = float64    // TODO: credibility
type k = struct{}   // TODO: k-order guid.K

// // index types for 3rd faction
// type __s = *skiplist.SkipList[s, k]
// type __p = *skiplist.SkipList[p, k]
type __o = *skiplist.SkipList[o, k]

// // index types for 2nd faction
type _po = *skiplist.SkipList[p, __o]

// type _op = *skiplist.SkipList[o, __p]
// type _so = *skiplist.SkipList[s, __o]
// type _os = *skiplist.SkipList[o, __s]
// type _sp = *skiplist.SkipList[s, __p]
// type _ps = *skiplist.SkipList[p, __s]

type spo = *skiplist.SkipList[s, _po]

type hspo = *skiplist.SkipList[h, spo]
