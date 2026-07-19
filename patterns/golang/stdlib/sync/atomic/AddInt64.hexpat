arch x86_64;
platform linux, darwin, windows, freebsd;

pattern atomic_AddInt64 { name: "sync/atomic.AddInt64"; library: "go-stdlib"; instr match { CALL sync_atomic_AddInt64 } gen { atomic.AddInt64($addr, $delta); } bind { addr as "addr"; delta as "delta"; } }
