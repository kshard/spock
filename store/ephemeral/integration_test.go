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

package ephemeral_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/fogfish/curie"
	"github.com/fogfish/it/v2"
	"github.com/kshard/spock"
	"github.com/kshard/spock/store/ephemeral"
)

const (
	A = curie.IRI("u:A")
	B = curie.IRI("u:B")
	C = curie.IRI("s:C")
	D = curie.IRI("u:D")
	E = curie.IRI("u:E")
	F = curie.IRI("s:F")
	G = curie.IRI("s:G")
	N = curie.IRI("n:N")
)

func datasetSocialGraph() spock.Bag {
	return spock.Bag{
		spock.From(A, "follows", B),
		spock.From(C, "follows", B),
		spock.From(C, "follows", E),
		spock.From(C, "relates", D),
		spock.From(D, "relates", B),
		spock.From(B, "follows", F),
		spock.From(F, "follows", G),
		spock.From(D, "relates", G),
		spock.From(E, "follows", F),

		spock.From(B, "status", "b"),
		spock.From(D, "status", "d"),
		spock.From(G, "status", "g"),
	}
}

func setup(bag spock.Bag) *ephemeral.Store {
	store := ephemeral.New()

	t := time.Now()
	ephemeral.Add(store, bag)

	fmt.Printf("==> setup %v\n", time.Since(t))

	return store
}

func TestSocialGraph(t *testing.T) {
	rds := setup(datasetSocialGraph())

	Seq := func(t *testing.T, uid string, req spock.Pattern) it.SeqOf[spock.SPOCK] {
		t.Helper()
		bag := spock.Bag{}
		seq, err := ephemeral.Match(rds, req)
		it.Then(t).Should(it.Nil(err))

		err = seq.FMap(bag.Join)
		it.Then(t).Should(
			it.Nil(err),
			it.Equal(req.String(), uid),
		)

		return it.Seq(bag)
	}

	//
	// #2: (s) ⇒ po
	//
	t.Run("#2: (s) ⇒ po", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(s) ⇒ po",
				spock.Query(spock.IRI.Equal(C), nil, nil),
			).Equal(
				spock.From(C, "follows", B),
				spock.From(C, "follows", E),
				spock.From(C, "relates", D),
			),
		)
	})

	t.Run("#2: (s) ⇒ po", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(s) ⇒ po",
				spock.Query(spock.IRI.Equal(N), nil, nil),
			).Equal(),
		)
	})

	//
	// #3: (sp) ⇒ o
	//
	t.Run("#3: (sp) ⇒ o", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(sp) ⇒ o",
				spock.Query(spock.IRI.Equal(C), spock.IRI.Equal("follows"), nil),
			).Equal(
				spock.From(C, "follows", B),
				spock.From(C, "follows", E),
			),
		)
	})

	t.Run("#3: (sp) ⇒ o", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(sp) ⇒ o",
				spock.Query(spock.IRI.Equal(C), spock.IRI.Equal("none"), nil),
			).Equal(),
		)
	})

	t.Run("#3: (sp) ⇒ o", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(sp) ⇒ o",
				spock.Query(spock.IRI.Equal(N), spock.IRI.Equal("follows"), nil),
			).Equal(),
		)
	})

	//
	// #4: (sᴾ) ⇒ o
	//

	t.Run("#4: (sᴾ) ⇒ o", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(sᴾ) ⇒ o",
				spock.Query(spock.IRI.Equal(C), spock.IRI.HasPrefix("f"), nil),
			).Equal(
				spock.From(C, "follows", B),
				spock.From(C, "follows", E),
			),
		)
	})

	t.Run("#4: (sᴾ) ⇒ o", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(sᴾ) ⇒ o",
				spock.Query(spock.IRI.Equal(C), spock.IRI.HasPrefix("n"), nil),
			).Equal(),
		)
	})

	t.Run("#4: (sᴾ) ⇒ o", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(sᴾ) ⇒ o",
				spock.Query(spock.IRI.Equal(N), spock.IRI.HasPrefix("f"), nil),
			).Equal(),
		)
	})

	//
	// #5: (so) ⇒ p
	//

	t.Run("#5: (so) ⇒ p", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(so) ⇒ p",
				spock.Query(spock.IRI.Equal(D), nil, spock.Eq(G)),
			).Equal(
				spock.From(D, "relates", G),
			),
		)
	})

	t.Run("#5: (so) ⇒ p", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(so) ⇒ p",
				spock.Query(spock.IRI.Equal(D), nil, spock.Eq("d")),
			).Equal(
				spock.From(D, "status", "d"),
			),
		)
	})

	t.Run("#5: (so) ⇒ p", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(so) ⇒ p",
				spock.Query(spock.IRI.Equal(D), nil, spock.Eq(N)),
			).Equal(),
		)
	})

	t.Run("#5: (so) ⇒ p", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(so) ⇒ p",
				spock.Query(spock.IRI.Equal(N), nil, spock.Eq(G)),
			).Equal(),
		)
	})

	//
	// #6: (sº) ⇒ p
	//

	t.Run("#6: (sº) ⇒ p", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(sº) ⇒ p",
				spock.Query(spock.IRI.Equal(D), nil, spock.HasPrefix(curie.IRI("s:"))),
			).Equal(
				spock.From(D, "relates", G),
			),
		)
	})

	t.Run("#6: (s)º ⇒ p", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(s)º ⇒ p",
				spock.Query(spock.IRI.Equal(D), nil, spock.Gt("a")),
			).Equal(
				spock.From(D, "status", "d"),
			),
		)
	})

	t.Run("#6: (s)º ⇒ p", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(s)º ⇒ p",
				spock.Query(spock.IRI.Equal(D), nil, spock.Lt("x")),
			).Equal(
				spock.From(D, "status", "d"),
			),
		)
	})

	t.Run("#6: (s)º ⇒ p", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(s)º ⇒ p",
				spock.Query(spock.IRI.Equal(D), nil, spock.Gt("x")),
			).Equal(),
		)
	})

	t.Run("#6: (s)º ⇒ p", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(s)º ⇒ p",
				spock.Query(spock.IRI.Equal(D), nil, spock.Lt("a")),
			).Equal(),
		)
	})

	//
	// #7: (spo) ⇒ ∅
	//

	t.Run("#7: (spo) ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(spo) ⇒ ∅",
				spock.Query(spock.IRI.Equal(C), spock.IRI.Equal("follows"), spock.Eq(E)),
			).Equal(
				spock.From(C, "follows", E),
			),
		)
	})

	t.Run("#7: (spo) ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(spo) ⇒ ∅",
				spock.Query(spock.IRI.Equal(C), spock.IRI.Equal("follows"), spock.Eq(N)),
			).Equal(),
		)
	})

	t.Run("#7: (spo) ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(spo) ⇒ ∅",
				spock.Query(spock.IRI.Equal(C), spock.IRI.Equal("none"), spock.Eq(N)),
			).Equal(),
		)
	})

	t.Run("#7: (spo) ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(spo) ⇒ ∅",
				spock.Query(spock.IRI.Equal(N), spock.IRI.Equal("none"), spock.Eq(N)),
			).Equal(),
		)
	})

	//
	// #8: (soᴾ) ⇒ ∅
	//

	t.Run("#8: (soᴾ) ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(soᴾ) ⇒ ∅",
				spock.Query(spock.IRI.Equal(C), spock.IRI.HasPrefix("f"), spock.Eq(E)),
			).Equal(
				spock.From(C, "follows", E),
			),
		)
	})

	t.Run("#8: (soᴾ) ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(soᴾ) ⇒ ∅",
				spock.Query(spock.IRI.Equal(C), spock.IRI.HasPrefix("n"), spock.Eq(E)),
			).Equal(),
		)
	})

	t.Run("#8: (so)ᴾ ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(soᴾ) ⇒ ∅",
				spock.Query(spock.IRI.Equal(C), spock.IRI.HasPrefix("f"), spock.Eq(N)),
			).Equal(),
		)
	})

	t.Run("#8: (so)ᴾ ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(soᴾ) ⇒ ∅",
				spock.Query(spock.IRI.Equal(N), spock.IRI.HasPrefix("f"), spock.Eq(E)),
			).Equal(),
		)
	})

	//
	// #9: (spº) ⇒ ∅
	//

	t.Run("#9: (spº) ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(spº) ⇒ ∅",
				spock.Query(spock.IRI.Equal(C), spock.IRI.Equal("follows"), spock.HasPrefix(curie.IRI("u:"))),
			).Equal(
				spock.From(C, "follows", B),
				spock.From(C, "follows", E),
			),
		)
	})

	t.Run("#9: (spº) ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(spº) ⇒ ∅",
				spock.Query(spock.IRI.Equal(C), spock.IRI.Equal("follows"), spock.HasPrefix(curie.IRI("n:"))),
			).Equal(),
		)
	})

	t.Run("#9: (sp)º ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(sp)º ⇒ ∅",
				spock.Query(spock.IRI.Equal(G), spock.IRI.Equal("status"), spock.Gt("a")),
			).Equal(
				spock.From(G, "status", "g"),
			),
		)
	})

	t.Run("#9: (sp)º ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(sp)º ⇒ ∅",
				spock.Query(spock.IRI.Equal(G), spock.IRI.Equal("status"), spock.Lt("x")),
			).Equal(
				spock.From(G, "status", "g"),
			),
		)
	})

	t.Run("#9: (sp)º ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(sp)º ⇒ ∅",
				spock.Query(spock.IRI.Equal(G), spock.IRI.Equal("status"), spock.Gt("x")),
			).Equal(),
		)
	})

	t.Run("#9: (sp)º ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(sp)º ⇒ ∅",
				spock.Query(spock.IRI.Equal(G), spock.IRI.Equal("status"), spock.Lt("a")),
			).Equal(),
		)
	})

	//
	// #10: (sᴾ)º ⇒ ∅
	//

	t.Run("#10: (sᴾº) ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(sᴾº) ⇒ ∅",
				spock.Query(spock.IRI.Equal(C), spock.IRI.HasPrefix("f"), spock.HasPrefix(curie.IRI("u:"))),
			).Equal(
				spock.From(C, "follows", B),
				spock.From(C, "follows", E),
			),
		)
	})

	t.Run("#10: (sᴾº) ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(sᴾº) ⇒ ∅",
				spock.Query(spock.IRI.Equal(C), spock.IRI.HasPrefix("f"), spock.HasPrefix(curie.IRI("n:"))),
			).Equal(),
		)
	})

	t.Run("#10: (sᴾº) ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(sᴾº) ⇒ ∅",
				spock.Query(spock.IRI.Equal(C), spock.IRI.HasPrefix("n"), spock.HasPrefix(curie.IRI("u:"))),
			).Equal(),
		)
	})

	t.Run("#10: (sᴾº) ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(sᴾº) ⇒ ∅",
				spock.Query(spock.IRI.Equal(N), spock.IRI.HasPrefix("f"), spock.HasPrefix(curie.IRI("u:"))),
			).Equal(),
		)
	})

	t.Run("#10: (sᴾ)º ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(sᴾ)º ⇒ ∅",
				spock.Query(spock.IRI.Equal(G), spock.IRI.HasPrefix("st"), spock.Gt("a")),
			).Equal(
				spock.From(G, "status", "g"),
			),
		)
	})

	t.Run("#10: (sᴾ)º ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(sᴾ)º ⇒ ∅",
				spock.Query(spock.IRI.Equal(G), spock.IRI.HasPrefix("st"), spock.Lt("x")),
			).Equal(
				spock.From(G, "status", "g"),
			),
		)
	})

	t.Run("#10: (sᴾ)º ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(sᴾ)º ⇒ ∅",
				spock.Query(spock.IRI.Equal(G), spock.IRI.HasPrefix("st"), spock.In("a", "x")),
			).Equal(
				spock.From(G, "status", "g"),
			),
		)
	})

	t.Run("#10: (sᴾ)º ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(sᴾ)º ⇒ ∅",
				spock.Query(spock.IRI.Equal(G), spock.IRI.HasPrefix("st"), spock.Gt("x")),
			).Equal(),
		)
	})

	t.Run("#10: (sᴾ)º ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(sᴾ)º ⇒ ∅",
				spock.Query(spock.IRI.Equal(G), spock.IRI.HasPrefix("st"), spock.Lt("a")),
			).Equal(),
		)
	})

	//
	// #11: (p) ⇒ so
	//

	t.Run("#11: (p) ⇒ so", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(p) ⇒ so",
				spock.Query(nil, spock.IRI.Equal("status"), nil),
			).Equal(
				spock.From(G, "status", "g"), // s:G < u:B
				spock.From(B, "status", "b"),
				spock.From(D, "status", "d"),
			),
		)
	})

	t.Run("#11: (p) ⇒ so", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(p) ⇒ so",
				spock.Query(nil, spock.IRI.Equal("none"), nil),
			).Equal(),
		)
	})

	//
	// #12: (po) ⇒ s
	//

	t.Run("#12: (po) ⇒ s", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(po) ⇒ s",
				spock.Query(nil, spock.IRI.Equal("follows"), spock.Eq(B)),
			).Equal(
				spock.From(C, "follows", B), // s:G < u:A
				spock.From(A, "follows", B),
			),
		)
	})
	t.Run("#12: (po) ⇒ s", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(po) ⇒ s",
				spock.Query(nil, spock.IRI.Equal("follows"), spock.Eq(N)),
			).Equal(),
		)
	})

	t.Run("#12: (po) ⇒ s", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(po) ⇒ s",
				spock.Query(nil, spock.IRI.Equal("none"), spock.Eq(B)),
			).Equal(),
		)
	})

	//
	// #13: (pº) ⇒ s
	//

	t.Run("#13: (pº) ⇒ s", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(pº) ⇒ s",
				spock.Query(nil, spock.IRI.Equal("follows"), spock.HasPrefix(curie.IRI("s:"))),
			).Equal(
				spock.From(B, "follows", F),
				spock.From(E, "follows", F),
				spock.From(F, "follows", G),
			),
		)
	})

	t.Run("#13: (p)º ⇒ s", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(p)º ⇒ s",
				spock.Query(nil, spock.IRI.Equal("status"), spock.Gt("a")),
			).Equal(
				spock.From(B, "status", "b"),
				spock.From(D, "status", "d"),
				spock.From(G, "status", "g"),
			),
		)
	})

	t.Run("#13: (p)º ⇒ s", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(p)º ⇒ s",
				spock.Query(nil, spock.IRI.Equal("status"), spock.Lt("x")),
			).Equal(
				spock.From(B, "status", "b"),
				spock.From(D, "status", "d"),
				spock.From(G, "status", "g"),
			),
		)
	})

	t.Run("#13: (p)º ⇒ s", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(p)º ⇒ s",
				spock.Query(nil, spock.IRI.Equal("status"), spock.In("d", "g")),
			).Equal(
				spock.From(D, "status", "d"),
				spock.From(G, "status", "g"),
			),
		)
	})

	t.Run("#13: (p)º ⇒ s", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(p)º ⇒ s",
				spock.Query(nil, spock.IRI.Equal("status"), spock.Gt("x")),
			).Equal(),
		)
	})

	t.Run("#13: (p)º ⇒ s", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(p)º ⇒ s",
				spock.Query(nil, spock.IRI.Equal("none"), spock.Gt("a")),
			).Equal(),
		)
	})

	//
	// #14: (pˢ) ⇒ o
	//

	t.Run("#14: (pˢ) ⇒ o", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(pˢ) ⇒ o",
				spock.Query(spock.IRI.HasPrefix("s:"), spock.IRI.Equal("follows"), nil),
			).Equal(
				spock.From(C, "follows", B),
				spock.From(C, "follows", E),
				spock.From(F, "follows", G),
			),
		)
	})

	t.Run("#14: (pˢ) ⇒ o", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(pˢ) ⇒ o",
				spock.Query(spock.IRI.HasPrefix("n:"), spock.IRI.Equal("follows"), nil),
			).Equal(),
		)
	})

	//
	// #15: (poˢ) ⇒ ∅
	//

	t.Run("#15: (poˢ) ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(poˢ) ⇒ ∅",
				spock.Query(spock.IRI.HasPrefix("s:"), spock.IRI.Equal("follows"), spock.Eq(E)),
			).Equal(
				spock.From(C, "follows", E),
			),
		)
	})

	t.Run("#15: (poˢ) ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(poˢ) ⇒ ∅",
				spock.Query(spock.IRI.HasPrefix("n:"), spock.IRI.Equal("follows"), spock.Eq(E)),
			).Equal(),
		)
	})

	t.Run("#15: (poˢ) ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(poˢ) ⇒ ∅",
				spock.Query(spock.IRI.HasPrefix("s:"), spock.IRI.Equal("follows"), spock.Eq(N)),
			).Equal(),
		)
	})

	t.Run("#15: (poˢ) ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(poˢ) ⇒ ∅",
				spock.Query(spock.IRI.HasPrefix("s:"), spock.IRI.Equal("none"), spock.Eq(E)),
			).Equal(),
		)
	})

	//
	// #16: (pˢ)º ⇒ ∅
	//

	t.Run("#16: (pˢº) ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(pˢº) ⇒ ∅",
				spock.Query(spock.IRI.HasPrefix("s:"), spock.IRI.Equal("follows"), spock.HasPrefix(curie.IRI("s:"))),
			).Equal(
				spock.From(F, "follows", G),
			),
		)
	})

	t.Run("#16: (pˢº) ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(pˢº) ⇒ ∅",
				spock.Query(spock.IRI.HasPrefix("s:"), spock.IRI.Equal("follows"), spock.HasPrefix(curie.IRI("n:"))),
			).Equal(),
		)
	})

	t.Run("#16: (pˢº) ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(pˢº) ⇒ ∅",
				spock.Query(spock.IRI.HasPrefix("n:"), spock.IRI.Equal("follows"), spock.HasPrefix(curie.IRI("s:"))),
			).Equal(),
		)
	})

	t.Run("#16: (pˢ)º ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(pˢ)º ⇒ ∅",
				spock.Query(spock.IRI.HasPrefix("s:"), spock.IRI.Equal("status"), spock.Gt("a")),
			).Equal(
				spock.From(G, "status", "g"),
			),
		)
	})

	t.Run("#16: (pˢ)º ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(pˢ)º ⇒ ∅",
				spock.Query(spock.IRI.HasPrefix("s:"), spock.IRI.Equal("status"), spock.Lt("x")),
			).Equal(
				spock.From(G, "status", "g"),
			),
		)
	})

	t.Run("#16: (pˢ)º ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(pˢ)º ⇒ ∅",
				spock.Query(spock.IRI.HasPrefix("s:"), spock.IRI.Equal("status"), spock.In("a", "x")),
			).Equal(
				spock.From(G, "status", "g"),
			),
		)
	})

	//
	// #17: (o) ⇒ ps
	//

	t.Run("#17: (o) ⇒ ps", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(o) ⇒ ps",
				spock.Query(nil, nil, spock.Eq(B)),
			).Equal(
				spock.From(C, "follows", B),
				spock.From(A, "follows", B),
				spock.From(D, "relates", B),
			),
		)
	})

	t.Run("#17: (o) ⇒ ps", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(o) ⇒ ps",
				spock.Query(nil, nil, spock.Eq(N)),
			).Equal(),
		)
	})

	//
	// #18: (oᴾ) ⇒ s
	//

	t.Run("#18: (oᴾ) ⇒ s", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(oᴾ) ⇒ s",
				spock.Query(nil, spock.IRI.HasPrefix("f"), spock.Eq(B)),
			).Equal(
				spock.From(C, "follows", B),
				spock.From(A, "follows", B),
			),
		)
	})

	t.Run("#18: (oᴾ) ⇒ s", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(oᴾ) ⇒ s",
				spock.Query(nil, spock.IRI.HasPrefix("n"), spock.Eq(B)),
			).Equal(),
		)
	})

	t.Run("#18: (oᴾ) ⇒ s", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(oᴾ) ⇒ s",
				spock.Query(nil, spock.IRI.HasPrefix("f"), spock.Eq(N)),
			).Equal(),
		)
	})

	//
	// #19: (oˢ) ⇒ p
	//

	t.Run("#19: (oˢ) ⇒ p", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(oˢ) ⇒ p",
				spock.Query(spock.IRI.HasPrefix("u:"), nil, spock.Eq(B)),
			).Equal(
				spock.From(A, "follows", B),
				spock.From(D, "relates", B),
			),
		)
	})

	t.Run("#19: (oˢ) ⇒ p", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(oˢ) ⇒ p",
				spock.Query(spock.IRI.HasPrefix("n:"), nil, spock.Eq(B)),
			).Equal(),
		)
	})

	t.Run("#19: (oˢ) ⇒ p", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(oˢ) ⇒ p",
				spock.Query(spock.IRI.HasPrefix("u:"), nil, spock.Eq(N)),
			).Equal(),
		)
	})

	//
	// #20: (oᴾˢ) ⇒ ∅
	//

	t.Run("#20: (oᴾˢ) ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(oᴾˢ) ⇒ ∅",
				spock.Query(spock.IRI.HasPrefix("u:"), spock.IRI.HasPrefix("f"), spock.Eq(B)),
			).Equal(
				spock.From(A, "follows", B),
			),
		)
	})

	t.Run("#20: (oᴾˢ) ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(oᴾˢ) ⇒ ∅",
				spock.Query(spock.IRI.HasPrefix("n:"), spock.IRI.HasPrefix("f"), spock.Eq(B)),
			).Equal(),
		)
	})

	t.Run("#20: (oᴾˢ) ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(oᴾˢ) ⇒ ∅",
				spock.Query(spock.IRI.HasPrefix("u:"), spock.IRI.HasPrefix("n"), spock.Eq(B)),
			).Equal(),
		)
	})

	//
	// #21: (ˢ) ⇒ po
	//

	t.Run("#21: (ˢ) ⇒ po", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(ˢ) ⇒ po",
				spock.Query(spock.IRI.HasPrefix("s:"), nil, nil),
			).Equal(
				spock.From(C, "follows", B),
				spock.From(C, "follows", E),
				spock.From(C, "relates", D),
				spock.From(F, "follows", G),
				spock.From(G, "status", "g"),
			),
		)
	})

	t.Run("#21: (ˢ) ⇒ po", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(ˢ) ⇒ po",
				spock.Query(spock.IRI.HasPrefix("n:"), nil, nil),
			).Equal(),
		)
	})

	//
	// #22: (ˢᴾ) ⇒ o
	//
	t.Run("#22: (ˢᴾ) ⇒ o", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(ˢᴾ) ⇒ o",
				spock.Query(spock.IRI.HasPrefix("s:"), spock.IRI.HasPrefix("f"), nil),
			).Equal(
				spock.From(C, "follows", B),
				spock.From(C, "follows", E),
				spock.From(F, "follows", G),
			),
		)
	})

	t.Run("#22: (ˢᴾ) ⇒ o", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(ˢᴾ) ⇒ o",
				spock.Query(spock.IRI.HasPrefix("s:"), spock.IRI.HasPrefix("n"), nil),
			).Equal(),
		)
	})

	t.Run("#22: (ˢᴾ) ⇒ o", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(ˢᴾ) ⇒ o",
				spock.Query(spock.IRI.HasPrefix("n:"), spock.IRI.HasPrefix("f"), nil),
			).Equal(),
		)
	})

	//
	// #23: (ˢº) ⇒ p
	//

	t.Run("#23: (ˢ)º ⇒ p", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(ˢ)º ⇒ p",
				spock.Query(spock.IRI.HasPrefix("s:"), nil, spock.Gt("a")),
			).Equal(
				spock.From(G, "status", "g"),
			),
		)
	})

	t.Run("#23: (ˢ)º ⇒ p", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(ˢ)º ⇒ p",
				spock.Query(spock.IRI.HasPrefix("s:"), nil, spock.Lt("x")),
			).Equal(
				spock.From(G, "status", "g"),
			),
		)
	})

	t.Run("#23: (ˢ)º ⇒ p", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(ˢ)º ⇒ p",
				spock.Query(spock.IRI.HasPrefix("s:"), nil, spock.Gt("x")),
			).Equal(),
		)
	})

	//
	// #24: (ˢᴾ)º ⇒ ∅
	//

	t.Run("#24: (ˢᴾ)º ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(ˢᴾ)º ⇒ ∅",
				spock.Query(spock.IRI.HasPrefix("s:"), spock.IRI.HasPrefix("s"), spock.Gt("a")),
			).Equal(
				spock.From(G, "status", "g"),
			),
		)
	})

	t.Run("#24: (ˢᴾ)º ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(ˢᴾ)º ⇒ ∅",
				spock.Query(spock.IRI.HasPrefix("s:"), spock.IRI.HasPrefix("s"), spock.Lt("x")),
			).Equal(
				spock.From(G, "status", "g"),
			),
		)
	})

	t.Run("#24: (ˢᴾ)º ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(ˢᴾ)º ⇒ ∅",
				spock.Query(spock.IRI.HasPrefix("s:"), spock.IRI.HasPrefix("s"), spock.Gt("x")),
			).Equal(),
		)
	})

	t.Run("#24: (ˢᴾ)º ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(ˢᴾ)º ⇒ ∅",
				spock.Query(spock.IRI.HasPrefix("s:"), spock.IRI.HasPrefix("n"), spock.Gt("a")),
			).Equal(),
		)
	})

	t.Run("#24: (ˢᴾ)º ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(ˢᴾ)º ⇒ ∅",
				spock.Query(spock.IRI.HasPrefix("n:"), spock.IRI.HasPrefix("s"), spock.Gt("a")),
			).Equal(),
		)
	})

	//
	// #25: (ᴾ) ⇒ so
	//

	t.Run("#25: (ᴾ) ⇒ so", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(ᴾ) ⇒ so",
				spock.Query(nil, spock.IRI.HasPrefix("rel"), nil),
			).Equal(
				spock.From(C, "relates", D),
				spock.From(D, "relates", G),
				spock.From(D, "relates", B),
			),
		)
	})

	t.Run("#25: (ᴾ) ⇒ so", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(ᴾ) ⇒ so",
				spock.Query(nil, spock.IRI.HasPrefix("n"), nil),
			).Equal(),
		)
	})

	//
	// #26: (ᴾ)º ⇒ s
	//

	t.Run("#26: (ᴾ)º ⇒ s", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(ᴾ)º ⇒ s",
				spock.Query(nil, spock.IRI.HasPrefix("s"), spock.Gt("a")),
			).Equal(
				spock.From(B, "status", "b"),
				spock.From(D, "status", "d"),
				spock.From(G, "status", "g"),
			),
		)
	})

	t.Run("#26: (ᴾ)º ⇒ s", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(ᴾ)º ⇒ s",
				spock.Query(nil, spock.IRI.HasPrefix("s"), spock.Lt("x")),
			).Equal(
				spock.From(B, "status", "b"),
				spock.From(D, "status", "d"),
				spock.From(G, "status", "g"),
			),
		)
	})

	t.Run("#26: (ᴾ)º ⇒ s", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(ᴾ)º ⇒ s",
				spock.Query(nil, spock.IRI.HasPrefix("s"), spock.In("c", "x")),
			).Equal(
				spock.From(D, "status", "d"),
				spock.From(G, "status", "g"),
			),
		)
	})

	t.Run("#26: (ᴾ)º ⇒ s", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(ᴾ)º ⇒ s",
				spock.Query(nil, spock.IRI.HasPrefix("s"), spock.Gt("x")),
			).Equal(),
		)
	})

	t.Run("#26: (ᴾ)º ⇒ s", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(ᴾ)º ⇒ s",
				spock.Query(nil, spock.IRI.HasPrefix("n"), spock.Gt("a")),
			).Equal(),
		)
	})

	//
	// #27: (º) ⇒ ps
	//

	t.Run("#27: (º) ⇒ ps", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(º) ⇒ ps",
				spock.Query(nil, nil, spock.HasPrefix(curie.IRI("u:"))),
			).Equal(
				spock.From(C, "follows", B),
				spock.From(A, "follows", B),
				spock.From(D, "relates", B),
				spock.From(C, "relates", D),
				spock.From(C, "follows", E),
			),
		)
	})

	t.Run("#27: ()º ⇒ ps", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "()º ⇒ ps",
				spock.Query(nil, nil, spock.Gt("a")),
			).Equal(
				spock.From(B, "status", "b"),
				spock.From(D, "status", "d"),
				spock.From(G, "status", "g"),
			),
		)
	})

	t.Run("#27: ()º ⇒ ps", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "()º ⇒ ps",
				spock.Query(nil, nil, spock.Lt("x")),
			).Equal(
				spock.From(B, "status", "b"),
				spock.From(D, "status", "d"),
				spock.From(G, "status", "g"),
			),
		)
	})

	t.Run("#27: ()º ⇒ ps", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "()º ⇒ ps",
				spock.Query(nil, nil, spock.In("c", "x")),
			).Equal(
				spock.From(D, "status", "d"),
				spock.From(G, "status", "g"),
			),
		)
	})

	t.Run("#27: ()º ⇒ ps", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "()º ⇒ ps",
				spock.Query(nil, nil, spock.Gt("x")),
			).Equal(),
		)
	})

}
