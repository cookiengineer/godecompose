arch x86_64;
platform linux, darwin, windows, freebsd;

pattern heap_Push { name: "container/heap.Push"; library: "go-stdlib"; instr match { CALL container_heap_Push } gen { heap.Push($h, $x); } bind { h as "h"; x as "x"; } }
