arch x86_64;
platform linux, darwin, windows, freebsd;

pattern ring_Move { name: "ring.Ring.Move"; library: "go-stdlib"; instr match { CALL container_ring___Ring__Move } gen { $r.Move($n); } bind { r as "r"; n as "n"; } }
