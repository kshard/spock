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

package jsonld

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
		graph, has := val["@graph"]
		if has {
			switch seq := graph.(type) {
			case []any:
				return decodeArray(bag, nil, nil, seq)
			default:
				return fmt.Errorf("json-ld graph codec do not support %T (%v)", val, val)
			}
		}
		return decodeObject(bag, nil, nil, val)
	default:
		return fmt.Errorf("json-ld codec do not support %T (%v)", val, val)
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
			return fmt.Errorf("json-ld array codec do not support %T (%v)", val, val)
		}
	}
	return nil
}

func decodeObject(bag *Bag, s, p *curie.IRI, obj map[string]any) error {
	uid, has := decodeObjectID(obj)
	if !has {
		uid = curie.New("_:%s", guid.L(guid.Clock))
	}

	if s != nil && p != nil {
		*bag = append(*bag, spock.From(*s, *p, uid))
	}

	isa, has := decodeObjectType(obj)
	if has {
		*bag = append(*bag, spock.From(uid, curie.IRI("rdf:type"), isa))
	}

	return decodeObjectProperties(bag, uid, obj)
}

func decodeObjectID(obj map[string]any) (curie.IRI, bool) {
	raw, has := obj["@id"]
	if !has {
		return "", false
	}

	id, ok := raw.(string)
	if !ok {
		return "", false
	}

	return curie.IRI(id), true
}

func decodeObjectType(obj map[string]any) (curie.IRI, bool) {
	if iri, has := decodeObjectTypeFrom("@type", obj); has {
		return iri, true
	}

	return decodeObjectTypeFrom("rdf:type", obj)
}

func decodeObjectTypeFrom(prop string, obj map[string]any) (curie.IRI, bool) {
	raw, has := obj[prop]
	if !has {
		return "", false
	}

	id, ok := raw.(string)
	if !ok {
		return "", false
	}

	return curie.IRI(id), true
}

func decodeObjectProperties(bag *Bag, s curie.IRI, obj map[string]any) error {
	for key, val := range obj {
		if key == "@id" || key == "@type" || key == "rdf:type" {
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
			if err := decodeNodeObject(bag, s, p, o); err != nil {
				return err
			}
		case []any:
			if err := decodeNodeArray(bag, s, p, o); err != nil {
				return err
			}
		default:
			return fmt.Errorf("json-ld object codec do not support %T (%v)", val, val)
		}
	}

	return nil
}

func decodeNodeObject(bag *Bag, s, p curie.IRI, node map[string]any) error {
	val, has := node["@value"]
	if has {
		return decodeValue(bag, s, p, val)
	}

	iri, has := decodeObjectID(node)
	if has {
		*bag = append(*bag, spock.From(s, p, iri))
		return nil
	}

	return fmt.Errorf("json-ld node object codec do not support %T (%v)", node, node)
}

func decodeNodeArray(bag *Bag, s, p curie.IRI, array []any) error {
	for _, val := range array {
		switch o := val.(type) {
		// case float64:
		// 	*bag = append(*bag, spock.From(s, p, o))
		case string:
			*bag = append(*bag, spock.From(s, p, o))
		// case bool:
		// 	*bag = append(*bag, spock.From(s, p, o))
		case map[string]any:
			decodeNodeObject(bag, s, p, o)
		default:
			return fmt.Errorf("json-ld node array codec do not support %T (%v)", val, val)
		}
	}

	return nil
}

func decodeValue(bag *Bag, s, p curie.IRI, val any) error {
	switch o := val.(type) {
	// case float64:
	// 	*bag = append(*bag, spock.From(s, p, o))
	case string:
		*bag = append(*bag, spock.From(s, p, o))
	// case bool:
	// 	*bag = append(*bag, spock.From(s, p, o))
	default:
		return fmt.Errorf("json-ld value codec do not support %T (%v)", val, val)
	}

	return nil
}
