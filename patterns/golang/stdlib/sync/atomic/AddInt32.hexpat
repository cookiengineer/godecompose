arch x86_64;
platform linux, darwin, windows, freebsd;

pattern atomic_AddInt32 { name: "sync/atomic.AddInt32"; library: "go-stdlib"; instr match { CALL sync_atomic_AddInt32 } gen { atomic.AddInt32($addr, $delta); } bind { addr as "addr"; delta as "delta"; } }
