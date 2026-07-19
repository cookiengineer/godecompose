arch x86_64;
platform linux, darwin, windows, freebsd;

pattern heap_Pop { name: "container/heap.Pop"; library: "go-stdlib"; instr match { CALL container_heap_Pop } gen { heap.Pop($h); } bind { h as "h"; } }
