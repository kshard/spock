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
// The file define DSL for predicate expressions
//

import (
	"fmt"

	"github.com/fogfish/curie"
	"github.com/kshard/xsd"
)

// types of predicate clauses
type Clause int

const (
	ALL Clause = iota
	EQ         // Equal
	PQ         // Prefix Equal
	LT         // Less Than
	GT         // Greater Than
	IN         // InRange, Between
)

// Predicate expression
type Predicate[T any] struct {
	Clause Clause
	Value  T
	Other  T
}

func (pred Predicate[T]) String() string {
	switch pred.Clause {
	case EQ:
		return fmt.Sprintf("= %v", pred.Value)
	case PQ:
		return fmt.Sprintf("~ %v", pred.Value)
	case LT:
		return fmt.Sprintf("< %v", pred.Value)
	case GT:
		return fmt.Sprintf("> %v", pred.Value)
	case IN:
		return fmt.Sprintf("[%v, %v]", pred.Value, pred.Other)
	default:
		return ""
	}
}

type iri string

const IRI = iri("")

// Makes `equal` to IRI predicate
func (iri) Eq(value curie.IRI) *Predicate[xsd.AnyURI] {
	return &Predicate[xsd.AnyURI]{Clause: EQ, Value: xsd.ToAnyURI(value)}
}

func (iri) Equal(value curie.IRI) *Predicate[xsd.AnyURI] {
	return IRI.Eq(value)
}

// Makes `prefix` to IRI predicate
// func (iri) HasPrefix(value curie.IRI) *Predicate[xsd.AnyURI] {
// 	return &Predicate[xsd.AnyURI]{Clause: PQ, Value: xsd.AnyURI(symbol.New(string(value)))}
// 	// return &Predicate[curie.IRI]{Clause: PQ, Value: value}
// }

// Makes `equal to` value predicate
// func Eq[T xsd.DataType](value T) *Predicate[xsd.Value] {
// 	return &Predicate[xsd.Value]{Clause: EQ, Value: xsd.From(value)}
// }

func Eq[T ~string](value T) *Predicate[xsd.Symbol] {
	return &Predicate[xsd.Symbol]{Clause: EQ, Value: xsd.ToSymbol(string(value))}
}

//

// Makes `prefix` value predicate
func HasPrefix[T xsd.DataType](value T) *Predicate[xsd.Value] {
	return &Predicate[xsd.Value]{Clause: PQ, Value: xsd.From(value)}
}

// Makes `less than` value predicate
func Lt[T xsd.DataType](value T) *Predicate[xsd.Value] {
	return &Predicate[xsd.Value]{Clause: LT, Value: xsd.From(value)}
}

// Makes `greater than` value predicate
func Gt[T xsd.DataType](value T) *Predicate[xsd.Value] {
	return &Predicate[xsd.Value]{Clause: GT, Value: xsd.From(value)}
}

// Makes `in range` predicate
func In[T xsd.DataType](from, to T) *Predicate[xsd.Value] {
	return &Predicate[xsd.Value]{Clause: IN, Value: xsd.From(from), Other: xsd.From(to)}
}
