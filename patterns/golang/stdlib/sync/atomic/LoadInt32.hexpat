arch x86_64;
platform linux, darwin, windows, freebsd;

pattern atomic_LoadInt32 { name: "sync/atomic.LoadInt32"; library: "go-stdlib"; instr match { CALL sync_atomic_LoadInt32 } gen { atomic.LoadInt32($addr); } bind { addr as "addr"; } }
