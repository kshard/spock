package main

import (
	"fmt"

	"github.com/kshard/spock/store/hybrid"
)

func main() {
	store := hybrid.New()
	hybrid.Config(store)

	// hybrid.Put(store, spock.From("a", "b", "1"))
	// hybrid.Put(store, spock.From("a", "b", "2"))
	// hybrid.Put(store, spock.From("a", "b", "3"))
	// hybrid.Put(store, spock.From("a", "d", "4"))
	// hybrid.Put(store, spock.From("a", "e", "5"))
	// hybrid.Put(store, spock.From("a", "f", "6"))

	s, err := hybrid.X(store)
	fmt.Println(err)

	for s.Next() {
		fmt.Println(s.Head())
	}
}
