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

package dynamo

import (
	"strings"

	"github.com/kshard/xsd"
)

//
// Pair codec
//

func encodeI(symbols Symbols, a xsd.AnyURI) string {
	prefix := symbols.ToSymbol(a).String()
	return prefix
}

func encodeII(symbols Symbols, a, b xsd.AnyURI) string {
	prefix := symbols.ToSymbol(a).String()
	suffix := symbols.ToSymbol(b).String()

	return prefix + "|" + suffix
}

func decodeII(symbols Symbols, val string) (xsd.AnyURI, xsd.AnyURI) {
	seq := strings.SplitN(val, "|", 2)
	prefix := symbols.FromString(seq[0])
	suffix := symbols.FromString(seq[1])

	return prefix, suffix
}

func encodeIV(symbols Symbols, a xsd.AnyURI, b xsd.Value) string {
	prefix := symbols.ToSymbol(a).String()
	suffix := encodeValue(symbols, b)

	return prefix + "|" + suffix
}

func decodeIV(symbols Symbols, val string) (xsd.AnyURI, xsd.Value) {
	seq := strings.SplitN(val, "|", 2)
	prefix := symbols.FromString(seq[0])
	suffix := decodeValue(symbols, seq[1])

	return prefix, suffix
}

func encodeVI(symbols Symbols, a xsd.Value, b xsd.AnyURI) string {
	prefix := encodeValue(symbols, a)
	suffix := symbols.ToSymbol(b).String()

	return prefix + "|" + suffix
}

func decodeVI(symbols Symbols, val string) (xsd.Value, xsd.AnyURI) {
	seq := strings.SplitN(val, "|", 2)
	prefix := decodeValue(symbols, seq[0])
	suffix := symbols.FromString(seq[1])

	return prefix, suffix

}

//
// Value codec - ᴸᴵᴳ
//

func encodeValue(symbols Symbols, value xsd.Value) string {
	switch v := value.(type) {
	case xsd.AnyURI:
		return "ᴵ" + symbols.ToSymbol(v).String()
	case xsd.String:
		return "ᴸ" + string(v)
	default:
		panic("not supported")
	}
}

func decodeValue(symbols Symbols, value string) xsd.Value {
	switch value[:3] {
	case "ᴵ":
		return symbols.FromString(value[3:])
	case "ᴸ":
		return xsd.String(value[3:])
	}

	return nil
}
