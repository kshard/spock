package hybrid

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/fogfish/curie"
	"github.com/fogfish/dynamo/v2/service/ddb"
	"github.com/fogfish/segment"
	"github.com/fogfish/skiplist"
)

type Dynamo[T skiplist.Num] struct {
	id  string
	idb *ddb.Storage[INode]
	xdb *ddb.Storage[Node]
}

func NewDynamo[T skiplist.Num](id string) *Dynamo[T] {
	idb, _ := ddb.New[INode]("ddb:///thing-test")
	xdb, _ := ddb.New[Node]("ddb:///thing-test")

	return &Dynamo[T]{
		id:  id,
		idb: idb,
		xdb: xdb,
	}
}

// type Segment struct {
// 	Rank uint32
// 	Lo   uint64
// 	Hi   uint64
// }

type INode struct {
	Prefix   string `dynamodbav:"prefix"`
	Suffix   string `dynamodbav:"suffix"`
	Capacity int    `dynamodbav:"capacity"`
	GF2      []byte `dynamodbav:"gf2"`
}

func (in INode) HashKey() curie.IRI { return curie.IRI(in.Prefix) }
func (in INode) SortKey() curie.IRI { return curie.IRI(in.Suffix) }

type Node struct {
	Prefix string `dynamodbav:"prefix"`
	Suffix string `dynamodbav:"suffix"`
	Data   []byte `dynamodbav:"data"`
}

func (in Node) HashKey() curie.IRI { return curie.IRI(in.Prefix) }
func (in Node) SortKey() curie.IRI { return curie.IRI(in.Suffix) }

func (f *Dynamo[T]) WriteINode(node *segment.INode[uint64]) error {
	pos := 0
	gf2 := make([]byte, node.GF2.Length*(4+8+8))
	seq := skiplist.ForGF2(node.GF2, node.GF2.Keys())
	for has := seq != nil; has; has = seq.Next() {
		arc := seq.Value()
		binary.BigEndian.PutUint32(gf2[pos+0x0:pos+0x4], arc.Rank)
		binary.BigEndian.PutUint64(gf2[pos+0x4:pos+0xc], arc.Lo)
		binary.BigEndian.PutUint64(gf2[pos+0xc:pos+0x14], arc.Hi)
		pos = pos + 0x14
	}

	inode := INode{
		Prefix:   "g",
		Suffix:   f.id,
		Capacity: node.Capacity,
		GF2:      gf2,
	}

	if err := f.idb.Put(context.Background(), inode); err != nil {
		return err
	}

	return nil
}

func (f *Dynamo[T]) Write(addr uint64, kv *skiplist.Map[uint64, *skiplist.Set[T]]) error {
	// TODO:
	//  - Symbols MUST BE globally unique due to segments topology

	if kv.Length == 0 {
		return nil
	}

	gen := map[uint64][]T{}
	seq := skiplist.ForMap(kv, kv.Keys())
	for has := seq != nil; has; has = seq.Next() {
		key := seq.Key()
		set := make([]T, 0)
		for e := seq.Value().Values(); e != nil; e = e.Next() {
			set = append(set, e.Key())
		}
		gen[key] = set
	}

	buf, _ := json.Marshal(gen)

	node := Node{
		Prefix: "g",
		Suffix: fmt.Sprintf("%s-%x", f.id, addr),
		Data:   buf,
	}

	if err := f.xdb.Put(context.Background(), node); err != nil {
		return err
	}

	return nil
}

func (f *Dynamo[T]) ReadINode(node *segment.INode[uint64]) error {
	inode, err := f.idb.Get(context.Background(), INode{Prefix: "g", Suffix: f.id})
	if err != nil {
		return nil
	}

	node.Capacity = inode.Capacity
	node.GF2 = skiplist.NewGF2[uint64]()

	for pos := 0; pos < len(inode.GF2); pos = pos + 0x14 {
		arc := skiplist.Arc[uint64]{}
		arc.Rank = binary.BigEndian.Uint32(inode.GF2[pos+0x0 : pos+0x4])
		arc.Lo = binary.BigEndian.Uint64(inode.GF2[pos+0x4 : pos+0xc])
		arc.Hi = binary.BigEndian.Uint64(inode.GF2[pos+0xc : pos+0x14])
		node.GF2.Put(arc)
	}

	return nil
}

func (f *Dynamo[T]) Read(addr uint64) (*skiplist.Map[uint64, *skiplist.Set[T]], error) {
	node, err := f.xdb.Get(context.Background(), Node{Prefix: "g", Suffix: fmt.Sprintf("%s-%x", f.id, addr)})
	if err != nil {
		fmt.Printf("==> err %v\n", err)
		return nil, nil
	}

	gen := map[uint64][]T{}
	if err := json.Unmarshal(node.Data, &gen); err != nil {
		fmt.Printf("==>> fuq %v\n", err)
		fmt.Println(string(node.Data))
	}

	kv := skiplist.NewMap[uint64, *skiplist.Set[T]]()
	for key, val := range gen {
		set := skiplist.NewSet[T]()
		for _, x := range val {
			set.Add(x)
		}
		kv.Put(key, set)
	}

	return kv, nil
}
