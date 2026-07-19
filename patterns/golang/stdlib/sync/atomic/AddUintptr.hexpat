arch x86_64;
platform linux, darwin, windows, freebsd;

pattern atomic_AddUintptr { name: "sync/atomic.AddUintptr"; library: "go-stdlib"; instr match { CALL sync_atomic_AddUintptr } gen { atomic.AddUintptr($addr, $delta); } bind { addr as "addr"; delta as "delta"; } }
