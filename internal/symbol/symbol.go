package symbol

import (
	"github.com/fogfish/curie"
	"github.com/fogfish/dynamo/v2"
	"github.com/fogfish/dynamo/v2/service/ddb"
	"github.com/fogfish/guid/v2"
)

// Globally unique identity for xsd.Symbol
type Symbol guid.K

func (sym Symbol) String() string { return guid.String(guid.K(sym)) }

// Persistent collection of symbols
type Symbols struct {
	Namespace curie.IRI `dynamodbav:"prefix"`
	ID        curie.IRI `dynamodbav:"suffix"`
	Symbols   []string  `dynamodbav:"symbols,stringset"`
}

func (seq Symbols) HashKey() curie.IRI { return seq.Namespace }
func (seq Symbols) SortKey() curie.IRI { return seq.ID }

func NewStore(connector string, opts ...dynamo.Option) (*ddb.Storage[*Symbols], error) {
	return ddb.New[*Symbols](connector, opts...)
}
