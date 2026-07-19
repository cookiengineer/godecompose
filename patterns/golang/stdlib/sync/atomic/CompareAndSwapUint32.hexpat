arch x86_64;
platform linux, darwin, windows, freebsd;

pattern atomic_CompareAndSwapUint32 { name: "sync/atomic.CompareAndSwapUint32"; library: "go-stdlib"; instr match { CALL sync_atomic_CompareAndSwapUint32 } gen { atomic.CompareAndSwapUint32($addr, $old, $new); } bind { addr as "addr"; old as "old"; new as "new"; } }
