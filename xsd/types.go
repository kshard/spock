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
	"strconv"

	"github.com/fogfish/curie"
)

type Value interface{ XSDType() curie.IRI }

// The data type represents Internationalized Resource Identifier.
// Used to uniquely identify concept, objects, etc.
type AnyURI curie.IRI

const XSD_ANYURI = curie.IRI("xsd:anyURI")

func (v AnyURI) XSDType() curie.IRI { return XSD_ANYURI }
func (v AnyURI) String() string     { return curie.IRI(v).Safe() }

// The string data-type represents character strings in knowledge statements.
// The language strings are annotated with corresponding language tag.
type String string

const XSD_STRING = curie.IRI("xsd:string")

func (v String) XSDType() curie.IRI { return XSD_STRING }
func (v String) String() string     { return strconv.Quote(string(v)) }

// The Integer data-type in knowledge statement.
// The library uses various int precision data-types to represent decimal values.
// type XSDInteger = int
// type Byte = int8
// type Short = int16
// type Int = int32
// type Long = int64
// type NonNegativeInteger = uint
// type UnsignedByte = uint8
// type UnsignedShort = uint16
// type UnsignedInt = uint32
// type UnsignedLong = uint64

// const XSD_INTEGER = curie.IRI("xsd:integer")

// type Integer struct{ Value XSDInteger }

// func (Integer) XSDType() curie.IRI { return XSD_INTEGER }
