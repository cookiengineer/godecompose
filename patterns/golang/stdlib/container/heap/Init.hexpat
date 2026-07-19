arch x86_64;
platform linux, darwin, windows, freebsd;

pattern heap_Init { name: "container/heap.Init"; library: "go-stdlib"; instr match { CALL container_heap_Init } gen { heap.Init($h); } bind { h as "h"; } }
