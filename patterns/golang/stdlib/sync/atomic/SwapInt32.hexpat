arch x86_64;
platform linux, darwin, windows, freebsd;

pattern atomic_SwapInt32 { name: "sync/atomic.SwapInt32"; library: "go-stdlib"; instr match { CALL sync_atomic_SwapInt32 } gen { atomic.SwapInt32($addr, $new); } bind { addr as "addr"; new as "new"; } }
