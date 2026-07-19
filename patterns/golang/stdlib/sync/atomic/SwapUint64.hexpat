arch x86_64;
platform linux, darwin, windows, freebsd;

pattern atomic_SwapUint64 { name: "sync/atomic.SwapUint64"; library: "go-stdlib"; instr match { CALL sync_atomic_SwapUint64 } gen { atomic.SwapUint64($addr, $new); } bind { addr as "addr"; new as "new"; } }
