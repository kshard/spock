package ephemeral

func Stats(store *Store) {
	// fmt.Printf("==> spo %d\n", skiplist.Length(store.spo))
	// fmt.Printf("==> sop %d\n", skiplist.Length(store.sop))
	// fmt.Printf("==> pso %d\n", skiplist.Length(store.pso))
	// fmt.Printf("==> pos %d\n", skiplist.Length(store.pos))
	// fmt.Printf("==> osp %d\n", skiplist.Length(store.osp))
	// fmt.Printf("==> ops %d\n", skiplist.Length(store.ops))

	// it := skiplist.Values(store.spo)
	// for it.Next() {
	// 	s, po := it.Head()
	// 	fmt.Printf("==> spo %v %d\n", s, skiplist.Length(po))
	// }

	// it := skiplist.Values(store.sop)
	// for it.Next() {
	// 	s, op := it.Head()
	// 	fmt.Printf("==> sop %v %d\n", s, skiplist.Length(op))
	// }

	// it := skiplist.Values(store.pso)
	// for it.Next() {
	// 	p, so := it.Head()
	// 	fmt.Printf("==> pso %v %d\n", p, skiplist.Length(so))

	// 	itt := skiplist.Values(so)
	// 	for itt.Next() {
	// 		s, o := itt.Head()
	// 		if skiplist.Length(o) > 4 {
	// 			fmt.Printf("==> pso %v %v, %d\n", p, s, skiplist.Length(o))
	// 		}
	// 	}
	// }

	// it := skiplist.Values(store.pos)
	// for it.Next() {
	// 	p, os := it.Head()
	// 	fmt.Printf("==> pos %v %d\n", p, skiplist.Length(os))

	// 	itt := skiplist.Values(os)
	// 	for itt.Next() {
	// 		o, s := itt.Head()
	// 		if skiplist.Length(s) > 4 {
	// 			fmt.Printf("==> pos %v %v, %d\n", p, o, skiplist.Length(s))
	// 		}
	// 	}
	// }

	// it := skiplist.Values(store.osp)
	// for it.Next() {
	// 	o, sp := it.Head()
	// 	fmt.Printf("==> osp %v %d\n", o, skiplist.Length(sp))
	// }

	// it := skiplist.Values(store.ops)
	// for it.Next() {
	// 	o, ps := it.Head()
	// 	fmt.Printf("==> ops %v %d\n", o, skiplist.Length(ps))
	// }
}
