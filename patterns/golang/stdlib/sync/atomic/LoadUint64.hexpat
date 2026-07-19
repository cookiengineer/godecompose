arch x86_64;
platform linux, darwin, windows, freebsd;

pattern atomic_LoadUint64 { name: "sync/atomic.LoadUint64"; library: "go-stdlib"; instr match { CALL sync_atomic_LoadUint64 } gen { atomic.LoadUint64($addr); } bind { addr as "addr"; } }
