arch x86_64;
platform linux, darwin, windows, freebsd;

pattern atomic_LoadUintptr { name: "sync/atomic.LoadUintptr"; library: "go-stdlib"; instr match { CALL sync_atomic_LoadUintptr } gen { atomic.LoadUintptr($addr); } bind { addr as "addr"; } }
