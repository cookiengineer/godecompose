arch x86_64;
platform linux, darwin, windows, freebsd;

pattern atomic_SwapUint32 { name: "sync/atomic.SwapUint32"; library: "go-stdlib"; instr match { CALL sync_atomic_SwapUint32 } gen { atomic.SwapUint32($addr, $new); } bind { addr as "addr"; new as "new"; } }
