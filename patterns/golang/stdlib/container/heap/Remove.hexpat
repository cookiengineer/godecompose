arch x86_64;
platform linux, darwin, windows, freebsd;

pattern heap_Remove { name: "container/heap.Remove"; library: "go-stdlib"; instr match { CALL container_heap_Remove } gen { heap.Remove($h, $i); } bind { h as "h"; i as "i"; } }
