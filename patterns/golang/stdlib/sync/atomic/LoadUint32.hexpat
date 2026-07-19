arch x86_64;
platform linux, darwin, windows, freebsd;

pattern atomic_LoadUint32 { name: "sync/atomic.LoadUint32"; library: "go-stdlib"; instr match { CALL sync_atomic_LoadUint32 } gen { atomic.LoadUint32($addr); } bind { addr as "addr"; } }
