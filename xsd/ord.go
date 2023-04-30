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
	"reflect"
	"strings"

	"github.com/fogfish/curie"
	"github.com/fogfish/skiplist/ord"
)

// HasPrefix tests whether the xsd value a begins with prefix b.
// It return false if notion of XSD type does not support prefix matching.
func HasPrefix(a, b Value) bool {
	switch av := a.(type) {
	case AnyURI:
		if bv, ok := b.(AnyURI); ok {
			return strings.HasPrefix(string(av), string(bv))
		}

		return false
	case String:
		if bv, ok := b.(String); ok {
			return strings.HasPrefix(string(av), string(bv))
		}

		return false
	}

	return false
}

func Compare(a, b Value) int {
	switch av := a.(type) {
	case AnyURI:
		if bv, ok := b.(AnyURI); ok {
			return compare(av, bv)
		}

		return compare(reflect.Kind(1000), typeOf(b))
	case String:
		if bv, ok := b.(String); ok {
			return compare(av, bv)
		}

		return compare(reflect.String, typeOf(b))
	}
	return 0
}

func typeOf(x any) reflect.Kind {
	switch x.(type) {
	case AnyURI:
		return reflect.Kind(1000)
	case String:
		return reflect.String
	// case string:
	// 	return reflect.String
	// case bool:
	// 	return reflect.Bool
	// case int:
	// 	return reflect.Int
	// case int8:
	// 	return reflect.Int8
	// case int16:
	// 	return reflect.Int16
	// case int32:
	// 	return reflect.Int32
	// case int64:
	// 	return reflect.Int64
	// case uint:
	// 	return reflect.Uint
	// case uint8:
	// 	return reflect.Uint8
	// case uint16:
	// 	return reflect.Uint16
	// case uint32:
	// 	return reflect.Uint32
	// case uint64:
	// 	return reflect.Uint64
	// case float32:
	// 	return reflect.Float32
	// case float64:
	// 	return reflect.Float64
	// case curie.IRI:
	// 	return reflect.Kind(1000)
	// case []byte:
	// 	return reflect.Kind(1001)
	default:
		panic("not supported")
		// return reflect.Invalid
	}
}

func compare[T interface {
	~string |
		~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}](a, b T) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}

const (
	OrdIRI = curieIRI("ord.iri")
	OrdXSD = xsdValue("ord.xsd")
)

// Type Class ord.Ord[curie.IRI]
type curieIRI string

func (curieIRI) Compare(a, b curie.IRI) int { return ord.String.Compare(string(a), string(b)) }

// Type Class ord.Ord[xsd.Value]
type xsdValue string

func (xsdValue) Compare(a, b Value) int { return Compare(a, b) }
