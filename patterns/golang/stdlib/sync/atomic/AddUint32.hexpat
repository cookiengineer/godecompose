arch x86_64;
platform linux, darwin, windows, freebsd;

pattern atomic_AddUint32 { name: "sync/atomic.AddUint32"; library: "go-stdlib"; instr match { CALL sync_atomic_AddUint32 } gen { atomic.AddUint32($addr, $delta); } bind { addr as "addr"; delta as "delta"; } }
