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
// The file define streaming protocol for hexastore
//

import (
	"github.com/kshard/xsd"
)

// Stream of knowledge statements ⟨s,p,o,c,k⟩
type Stream interface {
	Head() SPOCK
	Next() bool
	FMap(func(SPOCK) error) error
}

type filter struct {
	pred   func(SPOCK) bool
	stream Stream
}

func (filter *filter) Head() SPOCK {
	return filter.stream.Head()
}

func (filter *filter) Next() bool {
	for {
		if !filter.stream.Next() {
			return false
		}

		if filter.pred(filter.stream.Head()) {
			return true
		}
	}
}

func (filter *filter) FMap(f func(SPOCK) error) error {
	for filter.Next() {
		if err := f(filter.Head()); err != nil {
			return err
		}
	}
	return nil
}

func NewFilter(pred func(SPOCK) bool, stream Stream) Stream {
	return &filter{pred: pred, stream: stream}
}

func NewFilterO(hint Hint, q *Predicate[xsd.Value], stream Stream) Stream {
	switch hint {
	case HINT_MATCH:
		return NewFilter(
			func(spock SPOCK) bool { return xsd.Compare(spock.O, q.Value) == 0 },
			stream,
		)
	case HINT_FILTER_PREFIX:
		return NewFilter(
			func(spock SPOCK) bool { return xsd.HasPrefix(spock.O, q.Value) },
			stream,
		)
	case HINT_FILTER:
		switch q.Clause {
		case LT:
			return NewFilter(
				func(spock SPOCK) bool { return xsd.Compare(spock.O, q.Value) == -1 },
				stream,
			)
		case GT:
			return NewFilter(
				func(spock SPOCK) bool { return xsd.Compare(spock.O, q.Value) == 1 },
				stream,
			)
		case IN:
			return NewFilter(
				func(spock SPOCK) bool {
					return xsd.Compare(spock.O, q.Value) >= 0 && xsd.Compare(spock.O, q.Other) <= 0
				},
				stream,
			)
		}
	}

	return stream
}

func NewFilterP(hint Hint, q *Predicate[xsd.AnyURI], stream Stream) Stream {
	switch hint {
	case HINT_MATCH:
		return NewFilter(
			func(spock SPOCK) bool { return spock.P == q.Value },
			stream,
		)
		// case HINT_FILTER_PREFIX:
		// 	return NewFilter(
		// 		func(spock SPOCK) bool { return strings.HasPrefix(string(spock.P), string(q.Value)) },
		// 		stream,
		// 	)
	}

	return stream
}

func NewFilterS(hint Hint, q *Predicate[xsd.AnyURI], stream Stream) Stream {
	switch hint {
	case HINT_MATCH:
		return NewFilter(
			func(spock SPOCK) bool { return spock.S == q.Value },
			stream,
		)
		// case HINT_FILTER_PREFIX:
		// 	return NewFilter(
		// 		func(spock SPOCK) bool { return strings.HasPrefix(string(spock.S), string(q.Value)) },
		// 		stream,
		// 	)
	}

	return stream
}
