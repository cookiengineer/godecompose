arch x86_64;
platform linux, darwin, windows, freebsd;

pattern ring_Len { name: "ring.Ring.Len"; library: "go-stdlib"; instr match { CALL container_ring___Ring__Len } gen { $r.Len(); } bind { r as "r"; } }
