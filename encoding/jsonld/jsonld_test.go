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

package jsonld_test

import (
	"encoding/json"
	"testing"

	"github.com/fogfish/curie"
	"github.com/fogfish/guid/v2"
	"github.com/fogfish/it/v2"
	"github.com/kshard/spock"
	"github.com/kshard/spock/encoding/jsonld"
)

func TestJsonLdUnmarshal(t *testing.T) {
	guid.Clock = guid.NewClockMock()
	luid := curie.IRI("_:5...............")

	Codec := func(t *testing.T, input string) it.SeqOf[spock.SPOCK] {
		t.Helper()
		bag := jsonld.Bag{}
		err := json.Unmarshal([]byte(input), &bag)
		it.Then(t).Should(it.Nil(err))

		return it.Seq(bag)
	}

	t.Run("OnlyProperty", func(t *testing.T) {
		it.Then(t).Should(
			Codec(t, `{
				"prop": "title"
			}`).Equal(
				spock.From(luid, "prop", "title"),
			),
		)
	})

	t.Run("PropertyWithID", func(t *testing.T) {
		it.Then(t).Should(
			Codec(t, `{
					"@id": "id",
					"prop": "title"
				}`).Equal(
				spock.From("id", "prop", "title"),
			),
		)
	})

	// t.Run("PropertyInt", func(t *testing.T) {
	// 	it.Then(t).Should(
	// 		Codec(t, `{
	// 			"prop": 10
	// 		}`).Equal(
	// 			spock.From(luid, "prop", 10),
	// 		),
	// 	)
	// })

	// t.Run("PropertyFloat", func(t *testing.T) {
	// 	it.Then(t).Should(
	// 		Codec(t, `{
	// 			"prop": 10.0
	// 		}`).Equal(
	// 			spock.From(luid, "prop", 10.0),
	// 		),
	// 	)
	// })

	// t.Run("PropertyBool", func(t *testing.T) {
	// 	it.Then(t).Should(
	// 		Codec(t, `{
	// 			"prop": true
	// 		}`).Equal(
	// 			spock.From(luid, "prop", true),
	// 		),
	// 	)
	// })

	t.Run("PropertyArray", func(t *testing.T) {
		it.Then(t).Should(
			Codec(t, `{
				"prop": ["a", "b", "c"]
			}`).Equal(
				spock.From(luid, "prop", "a"),
				spock.From(luid, "prop", "b"),
				spock.From(luid, "prop", "c"),
			),
		)
	})

	// t.Run("PropertyArrayHeterogenous", func(t *testing.T) {
	// 	it.Then(t).Should(
	// 		Codec(t, `{
	// 			"prop": [1, "b", true]
	// 		}`).Equal(
	// 			spock.From(luid, "prop", 1),
	// 			spock.From(luid, "prop", "b"),
	// 			spock.From(luid, "prop", true),
	// 		),
	// 	)
	// })

	t.Run("ArrayOfObjects", func(t *testing.T) {
		it.Then(t).Should(
			Codec(t, `[
				{"@id": "id", "prop": "a"},
				{"@id": "id", "porp": "b"}
			]`).Equal(
				spock.From("id", "prop", "a"),
				spock.From("id", "porp", "b"),
			),
		)
	})

	t.Run("Graph", func(t *testing.T) {
		it.Then(t).Should(
			Codec(t, `{
				"@graph": [
					{
						"@id": "a",
						"prop": {"@id": "b"},
						"porp": {"@id": "c"}
					},
					{
						"@id": "b",
						"prop": {"@value": "title"}
					},
					{
						"@id": "c",
						"prop": {"@value": "title"}
					}
				]
			}`).Equal(
				spock.From("a", "prop", curie.IRI("b")),
				spock.From("a", "porp", curie.IRI("c")),
				spock.From("b", "prop", "title"),
				spock.From("c", "prop", "title"),
			),
		)
	})
}
