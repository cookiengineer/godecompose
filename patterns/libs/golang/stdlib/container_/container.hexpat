arch x86_64;
platform linux, darwin, windows, freebsd;

pattern list_New { name: "container/list.New"; library: "go-stdlib"; instr match { CALL container_list_New } gen { list.New(); } }
pattern list_Init { name: "list.List.Init"; library: "go-stdlib"; instr match { CALL container_list___List__Init } gen { $l.Init(); } bind { l as "l"; } }
pattern list_Len { name: "list.List.Len"; library: "go-stdlib"; instr match { CALL container_list___List__Len } gen { $l.Len(); } bind { l as "l"; } }
pattern list_Front { name: "list.List.Front"; library: "go-stdlib"; instr match { CALL container_list___List__Front } gen { $l.Front(); } bind { l as "l"; } }
pattern list_Back { name: "list.List.Back"; library: "go-stdlib"; instr match { CALL container_list___List__Back } gen { $l.Back(); } bind { l as "l"; } }
pattern list_PushFront { name: "list.List.PushFront"; library: "go-stdlib"; instr match { CALL container_list___List__PushFront } gen { $l.PushFront($v); } bind { l as "l"; v as "v"; } }
pattern list_PushBack { name: "list.List.PushBack"; library: "go-stdlib"; instr match { CALL container_list___List__PushBack } gen { $l.PushBack($v); } bind { l as "l"; v as "v"; } }
pattern list_InsertBefore { name: "list.List.InsertBefore"; library: "go-stdlib"; instr match { CALL container_list___List__InsertBefore } gen { $l.InsertBefore($v, $mark); } bind { l as "l"; v as "v"; mark as "mark"; } }
pattern list_InsertAfter { name: "list.List.InsertAfter"; library: "go-stdlib"; instr match { CALL container_list___List__InsertAfter } gen { $l.InsertAfter($v, $mark); } bind { l as "l"; v as "v"; mark as "mark"; } }
pattern list_MoveToFront { name: "list.List.MoveToFront"; library: "go-stdlib"; instr match { CALL container_list___List__MoveToFront } gen { $l.MoveToFront($e); } bind { l as "l"; e as "e"; } }
pattern list_MoveToBack { name: "list.List.MoveToBack"; library: "go-stdlib"; instr match { CALL container_list___List__MoveToBack } gen { $l.MoveToBack($e); } bind { l as "l"; e as "e"; } }
pattern list_Remove { name: "list.List.Remove"; library: "go-stdlib"; instr match { CALL container_list___List__Remove } gen { $l.Remove($e); } bind { l as "l"; e as "e"; } }

pattern heap_Init { name: "container/heap.Init"; library: "go-stdlib"; instr match { CALL container_heap_Init } gen { heap.Init($h); } bind { h as "h"; } }
pattern heap_Push { name: "container/heap.Push"; library: "go-stdlib"; instr match { CALL container_heap_Push } gen { heap.Push($h, $x); } bind { h as "h"; x as "x"; } }
pattern heap_Pop { name: "container/heap.Pop"; library: "go-stdlib"; instr match { CALL container_heap_Pop } gen { heap.Pop($h); } bind { h as "h"; } }
pattern heap_Remove { name: "container/heap.Remove"; library: "go-stdlib"; instr match { CALL container_heap_Remove } gen { heap.Remove($h, $i); } bind { h as "h"; i as "i"; } }
pattern heap_Fix { name: "container/heap.Fix"; library: "go-stdlib"; instr match { CALL container_heap_Fix } gen { heap.Fix($h, $i); } bind { h as "h"; i as "i"; } }

pattern ring_New { name: "container/ring.New"; library: "go-stdlib"; instr match { CALL container_ring_New } gen { ring.New($n); } bind { n as "n"; } }
pattern ring_Len { name: "ring.Ring.Len"; library: "go-stdlib"; instr match { CALL container_ring___Ring__Len } gen { $r.Len(); } bind { r as "r"; } }
pattern ring_Next { name: "ring.Ring.Next"; library: "go-stdlib"; instr match { CALL container_ring___Ring__Next } gen { $r.Next(); } bind { r as "r"; } }
pattern ring_Prev { name: "ring.Ring.Prev"; library: "go-stdlib"; instr match { CALL container_ring___Ring__Prev } gen { $r.Prev(); } bind { r as "r"; } }
pattern ring_Link { name: "ring.Ring.Link"; library: "go-stdlib"; instr match { CALL container_ring___Ring__Link } gen { $r.Link($s); } bind { r as "r"; s as "s"; } }
pattern ring_Unlink { name: "ring.Ring.Unlink"; library: "go-stdlib"; instr match { CALL container_ring___Ring__Unlink } gen { $r.Unlink($n); } bind { r as "r"; n as "n"; } }
pattern ring_Move { name: "ring.Ring.Move"; library: "go-stdlib"; instr match { CALL container_ring___Ring__Move } gen { $r.Move($n); } bind { r as "r"; n as "n"; } }
pattern ring_Do { name: "ring.Ring.Do"; library: "go-stdlib"; instr match { CALL container_ring___Ring__Do } gen { $r.Do($f); } bind { r as "r"; f as "f"; } }
