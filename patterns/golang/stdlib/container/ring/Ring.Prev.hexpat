arch x86_64;
platform linux, darwin, windows, freebsd;

pattern ring_Prev { name: "ring.Ring.Prev"; library: "go-stdlib"; instr match { CALL container_ring___Ring__Prev } gen { $r.Prev(); } bind { r as "r"; } }
