arch x86_64;
platform linux, darwin, windows, freebsd;

pattern ring_Unlink { name: "ring.Ring.Unlink"; library: "go-stdlib"; instr match { CALL container_ring___Ring__Unlink } gen { $r.Unlink($n); } bind { r as "r"; n as "n"; } }
