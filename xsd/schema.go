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

package xsd

import (
	"fmt"

	"github.com/fogfish/curie"
)

// DataType is a type constrain used by the library.
// See https://www.w3.org/TR/xmlschema-2/#datatype
//
// Knowledge statements contain scalar objects -- literals. Literals are either
// language-tagged string `rdf:langString` or type-safe values containing a
// reference to data-type (e.g. `xsd:string`).
//
// This interface defines data-types supported by the library. It maps well-known
// semantic types to Golang native types and relation to existed schema(s) and
// ontologies.
type DataType interface {
	~string |
		~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64 |
		~bool |
		~[]byte
}

// The floating point data-type in knowledge statement.
// The library uses various uint precisions.
type Float = float32
type Double = float64

// The boolean data-type in knowledge statement
type Boolean = bool

type HexBinary = []byte
type Base64Binary = []byte

// From builds Object from Golang type
func From[T DataType](value T) Value {
	switch v := any(value).(type) {
	case curie.IRI:
		return AnyURI(v)
	case AnyURI:
		return v
	case string:
		return String(v)
	case String:
		return v
	// case int:
	// 	return Integer{Value: v}
	default:
		panic(fmt.Errorf("package xsd does not support %T", value))
	}
}
