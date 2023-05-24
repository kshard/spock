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
// The file define knowledge statement type
//

import (
	"fmt"

	"github.com/fogfish/curie"
	"github.com/fogfish/guid/v2"
	"github.com/kshard/xsd"
)

// Knowledge statement
type SPOCK struct {
	S xsd.AnyURI // s: subject
	P xsd.AnyURI // p: predicate
	O xsd.Symbol // o: object
	C float64    // c: credibility
	K guid.K     // k: k-order
}

func (spock SPOCK) String() string {
	return fmt.Sprintf("⟨%s %s %s⟩", spock.S, spock.P, spock.O)
}

// Create new knowledge statement From
func From[T ~string](s, p curie.IRI, o T) SPOCK {
	return SPOCK{S: xsd.ToAnyURI(s), P: xsd.ToAnyURI(p), O: xsd.ToSymbol(string(o))}
}

// func From(s, p curie.IRI, o xsd.Symbol) SPOCK {
// 	return SPOCK{S: xsd.ToAnyURI(s), P: xsd.ToAnyURI(p), O: o}
// }

// Collection of knowledge statements
type Bag []SPOCK

func (bag *Bag) Join(spock SPOCK) error {
	*bag = append(*bag, spock)
	return nil
}
