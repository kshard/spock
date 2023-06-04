package gob

import (
	"bytes"
	"encoding/gob"
	"io"

	"github.com/kshard/spock"
	"github.com/kshard/xsd"
)

func init() {
	gob.Register(xsd.AnyURI(0))
	gob.Register(xsd.String(""))
}

func Encode(bag spock.Bag) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(bag); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Decode(reader io.Reader) (spock.Bag, error) {
	var bag spock.Bag
	dec := gob.NewDecoder(reader)
	if err := dec.Decode(&bag); err != nil {
		return nil, err
	}

	return bag, nil
}
