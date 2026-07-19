arch x86_64;
platform linux, darwin, windows, freebsd;

pattern heap_Fix { name: "container/heap.Fix"; library: "go-stdlib"; instr match { CALL container_heap_Fix } gen { heap.Fix($h, $i); } bind { h as "h"; i as "i"; } }
