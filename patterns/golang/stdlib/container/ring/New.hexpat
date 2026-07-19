arch x86_64;
platform linux, darwin, windows, freebsd;

pattern ring_New { name: "container/ring.New"; library: "go-stdlib"; instr match { CALL container_ring_New } gen { ring.New($n); } bind { n as "n"; } }
