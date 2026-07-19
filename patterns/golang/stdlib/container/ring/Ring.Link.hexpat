arch x86_64;
platform linux, darwin, windows, freebsd;

pattern ring_Link { name: "ring.Ring.Link"; library: "go-stdlib"; instr match { CALL container_ring___Ring__Link } gen { $r.Link($s); } bind { r as "r"; s as "s"; } }
