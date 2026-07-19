arch x86_64;
platform linux, darwin, windows, freebsd;

pattern atomic_AddUint64 { name: "sync/atomic.AddUint64"; library: "go-stdlib"; instr match { CALL sync_atomic_AddUint64 } gen { atomic.AddUint64($addr, $delta); } bind { addr as "addr"; delta as "delta"; } }
