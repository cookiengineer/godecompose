arch x86_64;
platform linux, darwin, windows, freebsd;

pattern ring_Next { name: "ring.Ring.Next"; library: "go-stdlib"; instr match { CALL container_ring___Ring__Next } gen { $r.Next(); } bind { r as "r"; } }
