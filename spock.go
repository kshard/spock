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
	"github.com/kshard/spock/xsd"
)

// Knowledge statement
//
//	s: subject
//	p: predicate
//	o: object
//	c: credibility
//	k: k-order
type SPOCK struct {
	S curie.IRI
	P curie.IRI
	O xsd.Value
	C float64
	K guid.K
}

func (spock SPOCK) String() string {
	return fmt.Sprintf("⟨%s %s %s⟩", spock.S.Safe(), spock.P.Safe(), spock.O)
}

// Create new knowledge statement From
func From[T xsd.DataType](s, p curie.IRI, o T) SPOCK {
	return SPOCK{S: s, P: p, O: xsd.From(o)}
}

// Collection of knowledge statements
type Bag []SPOCK

func (bag *Bag) Join(spock SPOCK) error {
	*bag = append(*bag, spock)
	return nil
}
