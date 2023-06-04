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

package json

import (
	"encoding/json"
	"fmt"

	"github.com/fogfish/curie"
	"github.com/fogfish/guid/v2"
	"github.com/kshard/spock"
)

type Bag spock.Bag

func (bag *Bag) UnmarshalJSON(b []byte) error {
	var raw any

	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	switch val := raw.(type) {
	case []any:
		return decodeArray(bag, nil, nil, val)
	case map[string]any:
		return decodeObject(bag, nil, nil, val)
	default:
		return fmt.Errorf("json codec do not support %T (%v)", val, val)
	}
}

func decodeArray(bag *Bag, s, p *curie.IRI, seq []any) error {
	for _, val := range seq {
		switch o := val.(type) {
		// case float64:
		// 	if s != nil && p != nil {
		// 		*bag = append(*bag, spock.From(*s, *p, o))
		// 	}
		case string:
			if s != nil && p != nil {
				*bag = append(*bag, spock.From(*s, *p, o))
			}
		// case bool:
		// 	if s != nil && p != nil {
		// 		*bag = append(*bag, spock.From(*s, *p, o))
		// 	}
		case map[string]any:
			decodeObject(bag, s, p, o)
		default:
			return fmt.Errorf("json array codec do not support %T (%v)", val, val)
		}
	}
	return nil
}

func decodeObject(bag *Bag, s, p *curie.IRI, obj map[string]any) error {
	id, has := decodeObjectID(obj)
	if !has {
		id = curie.New("_:%s", guid.L(guid.Clock))
	}

	if s != nil && p != nil {
		*bag = append(*bag, spock.From(*s, *p, id))
	}

	return decodeObjectProperties(bag, id, obj)
}

func decodeObjectID(obj map[string]any) (curie.IRI, bool) {
	raw, has := obj["@id"]
	if !has {
		raw, has = obj["id"]
		if !has {
			return "", false
		}
	}

	id, ok := raw.(string)
	if !ok {
		return "", false
	}

	return curie.IRI(id), true
}

func decodeObjectProperties(bag *Bag, s curie.IRI, obj map[string]any) error {
	for key, val := range obj {
		if key == "@id" || key == "id" {
			continue
		}
		p := curie.IRI(key)

		switch o := val.(type) {
		// case float64:
		// 	*bag = append(*bag, spock.From(s, p, o))
		case string:
			*bag = append(*bag, spock.From(s, p, o))
		// case bool:
		// 	*bag = append(*bag, spock.From(s, p, o))
		case map[string]any:
			if err := decodeObject(bag, &s, &p, o); err != nil {
				return err
			}
		case []any:
			if err := decodeArray(bag, &s, &p, o); err != nil {
				return err
			}
		default:
			return fmt.Errorf("json object codec do not support %T (%v)", val, val)
		}
	}

	return nil
}
