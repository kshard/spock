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

package spock

//
// The file define query planner for hexastore
//

import (
	"fmt"

	"github.com/fogfish/curie"
	"github.com/kshard/spock/xsd"
)

// Lookup strategy for knowledge statements
type Strategy int

const (
	STRATEGY_NONE Strategy = iota
	STRATEGY_SPO
	STRATEGY_SOP
	STRATEGY_PSO
	STRATEGY_POS
	STRATEGY_OPS
	STRATEGY_OSP
)

// Hint to construct the query
type Hint int

const (
	HINT_NONE Hint = iota
	HINT_FILTER
	HINT_FILTER_PREFIX
	HINT_MATCH
)

type Pattern struct {
	Strategy                     Strategy
	S                            *Predicate[curie.IRI]
	P                            *Predicate[curie.IRI]
	O                            *Predicate[xsd.Value]
	HintForS, HintForP, HintForO Hint
}

func (q Pattern) toStringS() (string, string, string) {
	switch q.HintForS {
	case HINT_MATCH:
		return "s", "", ""
	case HINT_FILTER_PREFIX:
		return "ˢ", "", ""
	case HINT_FILTER:
		return "", "ˢ", ""
	default:
		return "", "", "s"
	}
}

func (q Pattern) toStringP() (string, string, string) {
	switch q.HintForP {
	case HINT_MATCH:
		return "p", "", ""
	case HINT_FILTER_PREFIX:
		return "ᴾ", "", ""
	case HINT_FILTER:
		return "", "ᴾ", ""
	default:
		return "", "", "p"
	}
}

func (q Pattern) toStringO() (string, string, string) {
	switch q.HintForO {
	case HINT_MATCH:
		return "o", "", ""
	case HINT_FILTER_PREFIX:
		return "º", "", ""
	case HINT_FILTER:
		return "", "º", ""
	default:
		return "", "", "o"
	}
}

func (q Pattern) String() string {
	s0, s1, s2 := q.toStringS()
	p0, p1, p2 := q.toStringP()
	o0, o1, o2 := q.toStringO()
	t := ""
	if s2 == "" && p2 == "" && o2 == "" {
		t = "∅"
	}

	switch {
	case q.Strategy == STRATEGY_SPO:
		return "(" + s0 + p0 + o0 + ")" + s1 + p1 + o1 + " ⇒ " + s2 + p2 + o2 + t
	case q.Strategy == STRATEGY_SOP:
		return "(" + s0 + o0 + p0 + ")" + s1 + o1 + p1 + " ⇒ " + s2 + o2 + p2 + t
	case q.Strategy == STRATEGY_PSO:
		return "(" + p0 + s0 + o0 + ")" + p1 + s1 + o1 + " ⇒ " + p2 + s2 + o2 + t
	case q.Strategy == STRATEGY_POS:
		return "(" + p0 + o0 + s0 + ")" + p1 + o1 + s1 + " ⇒ " + p2 + o2 + s2 + t
	case q.Strategy == STRATEGY_OPS:
		return "(" + o0 + p0 + s0 + ")" + o1 + p1 + s1 + " ⇒ " + o2 + p2 + s2 + t
	case q.Strategy == STRATEGY_OSP:
		return "(" + o0 + s0 + p0 + ")" + o1 + s1 + p1 + " ⇒ " + o2 + s2 + p2 + t
	}

	return "(___) ⇒ ∅"
}

func (q Pattern) Dump() string {
	return fmt.Sprintf("⟪%s : s %s, p %s, o %s⟫", q.String(), q.S, q.P, q.O)
}

func Query(
	s *Predicate[curie.IRI],
	p *Predicate[curie.IRI],
	o *Predicate[xsd.Value],
) Pattern {
	q := Pattern{
		S: s, P: p, O: o,
		HintForS: hintFor(s),
		HintForP: hintFor(p),
		HintForO: hintFor(o),
	}
	q.Strategy = strategy(q)

	return q
}

func hintFor[T any](pred *Predicate[T]) Hint {
	switch {
	case pred != nil && pred.Clause != EQ && pred.Clause != PQ:
		return HINT_FILTER
	case pred != nil && pred.Clause == PQ:
		return HINT_FILTER_PREFIX
	case pred != nil && pred.Clause == EQ:
		return HINT_MATCH
	default:
		return HINT_NONE
	}
}

// Estimates execution strategy for pattern
func strategy(q Pattern) Strategy {
	switch {
	case q.HintForS == HINT_MATCH:
		return strategyForS(q)
	case q.HintForP == HINT_MATCH:
		return strategyForP(q)
	case q.HintForO == HINT_MATCH:
		return strategyForO(q)
	// #1: ___ ⇒ spo
	case q.HintForS == HINT_NONE && q.HintForP == HINT_NONE && q.HintForO == HINT_NONE:
		return STRATEGY_NONE
	default:
		return strategyForX(q)
	}
}

func strategyForS(q Pattern) Strategy {
	switch {
	// #2: x__ ⇒ spo
	case q.HintForP == HINT_NONE && q.HintForO == HINT_NONE:
		return STRATEGY_SPO
	// #3: xx_ ⇒ spo
	case q.HintForP == HINT_MATCH && q.HintForO == HINT_NONE:
		return STRATEGY_SPO
	// #4: xo_ ⇒ spo
	case (q.HintForP == HINT_FILTER_PREFIX || q.HintForP == HINT_FILTER) && q.HintForO == HINT_NONE:
		return STRATEGY_SPO
	// #5: x_x ⇒ sop
	case q.HintForP == HINT_NONE && q.HintForO == HINT_MATCH:
		return STRATEGY_SOP
	// #6: x_o ⇒ sop
	case q.HintForP == HINT_NONE && (q.HintForO == HINT_FILTER_PREFIX || q.HintForO == HINT_FILTER):
		return STRATEGY_SOP
	// #7: xxx ⇒ spo
	case q.HintForP == HINT_MATCH && q.HintForO == HINT_MATCH:
		return STRATEGY_SPO
	// #8: xox ⇒ sop
	case (q.HintForP == HINT_FILTER_PREFIX || q.HintForP == HINT_FILTER) && q.HintForO == HINT_MATCH:
		return STRATEGY_SOP
	// #9: xxo ⇒ spo
	case q.HintForP == HINT_MATCH && (q.HintForO == HINT_FILTER_PREFIX || q.HintForO == HINT_FILTER):
		return STRATEGY_SPO
	// #10: xoo ⇒ spo
	case (q.HintForP == HINT_FILTER_PREFIX || q.HintForP == HINT_FILTER) && (q.HintForO == HINT_FILTER_PREFIX || q.HintForO == HINT_FILTER):
		return STRATEGY_SPO
	default:
		return STRATEGY_NONE
	}
}

func strategyForP(q Pattern) Strategy {
	switch {
	// #11: _x_ ⇒ pso
	case q.HintForS == HINT_NONE && q.HintForO == HINT_NONE:
		return STRATEGY_PSO
	// #12: _xx ⇒ pos
	case q.HintForS == HINT_NONE && q.HintForO == HINT_MATCH:
		return STRATEGY_POS
	// #13: _xo ⇒ pos
	case q.HintForS == HINT_NONE && (q.HintForO == HINT_FILTER_PREFIX || q.HintForO == HINT_FILTER):
		return STRATEGY_POS
	// #14: ox_ ⇒ pso
	case (q.HintForS == HINT_FILTER_PREFIX || q.HintForS == HINT_FILTER) && q.HintForO == HINT_NONE:
		return STRATEGY_PSO
	// #15: oxx ⇒ pos
	case (q.HintForS == HINT_FILTER_PREFIX || q.HintForS == HINT_FILTER) && q.HintForO == HINT_MATCH:
		return STRATEGY_POS
	// #16: oxo ⇒ pso
	case (q.HintForS == HINT_FILTER_PREFIX || q.HintForS == HINT_FILTER) && (q.HintForO == HINT_FILTER_PREFIX || q.HintForO == HINT_FILTER):
		return STRATEGY_PSO
	default:
		return STRATEGY_NONE
	}
}

func strategyForO(q Pattern) Strategy {
	switch {
	// #17: __x ⇒ ops
	case q.HintForS == HINT_NONE && q.HintForP == HINT_NONE:
		return STRATEGY_OPS
	// #18: _ox ⇒ ops
	case q.HintForS == HINT_NONE && (q.HintForP == HINT_FILTER_PREFIX || q.HintForP == HINT_FILTER):
		return STRATEGY_OPS
	// #19: o_x ⇒ osp
	case (q.HintForS == HINT_FILTER_PREFIX || q.HintForS == HINT_FILTER) && q.HintForP == HINT_NONE:
		return STRATEGY_OSP
	// #20: oox ⇒ ops
	case (q.HintForS == HINT_FILTER_PREFIX || q.HintForS == HINT_FILTER) && (q.HintForP == HINT_FILTER_PREFIX || q.HintForP == HINT_FILTER):
		return STRATEGY_OPS
	default:
		return STRATEGY_NONE
	}
}

func strategyForX(q Pattern) Strategy {
	switch {
	// #21: o__ ⇒ spo
	case (q.HintForS == HINT_FILTER_PREFIX || q.HintForS == HINT_FILTER) && q.HintForP == HINT_NONE && q.HintForO == HINT_NONE:
		return STRATEGY_SPO
	// #22: oo_ ⇒ spo
	case (q.HintForS == HINT_FILTER_PREFIX || q.HintForS == HINT_FILTER) && (q.HintForP == HINT_FILTER_PREFIX || q.HintForP == HINT_FILTER) && q.HintForO == HINT_NONE:
		return STRATEGY_SPO
	// #23: o_o ⇒ sop
	case (q.HintForS == HINT_FILTER_PREFIX || q.HintForS == HINT_FILTER) && q.HintForP == HINT_NONE && (q.HintForO == HINT_FILTER_PREFIX || q.HintForO == HINT_FILTER):
		return STRATEGY_SOP
	// #24: ooo  ⇒ spo
	case (q.HintForS == HINT_FILTER_PREFIX || q.HintForS == HINT_FILTER) && (q.HintForP == HINT_FILTER_PREFIX || q.HintForP == HINT_FILTER) && (q.HintForO == HINT_FILTER_PREFIX || q.HintForO == HINT_FILTER):
		return STRATEGY_SPO
	// #25: _o_ ⇒ pso
	case q.HintForS == HINT_NONE && (q.HintForP == HINT_FILTER_PREFIX || q.HintForP == HINT_FILTER) && q.HintForO == HINT_NONE:
		return STRATEGY_PSO
	// #26: _oo ⇒ pos
	case q.HintForS == HINT_NONE && (q.HintForP == HINT_FILTER_PREFIX || q.HintForP == HINT_FILTER) && (q.HintForO == HINT_FILTER_PREFIX || q.HintForO == HINT_FILTER):
		return STRATEGY_POS
	// #27: __o ⇒ ops
	case q.HintForS == HINT_NONE && q.HintForP == HINT_NONE && (q.HintForO == HINT_FILTER_PREFIX || q.HintForO == HINT_FILTER):
		return STRATEGY_OPS
	default:
		return STRATEGY_NONE
	}
}
