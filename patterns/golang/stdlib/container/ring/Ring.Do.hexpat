arch x86_64;
platform linux, darwin, windows, freebsd;

pattern ring_Do { name: "ring.Ring.Do"; library: "go-stdlib"; instr match { CALL container_ring___Ring__Do } gen { $r.Do($f); } bind { r as "r"; f as "f"; } }
