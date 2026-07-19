arch x86_64;
platform linux, darwin, windows, freebsd;

pattern atomic_CompareAndSwapUint64 { name: "sync/atomic.CompareAndSwapUint64"; library: "go-stdlib"; instr match { CALL sync_atomic_CompareAndSwapUint64 } gen { atomic.CompareAndSwapUint64($addr, $old, $new); } bind { addr as "addr"; old as "old"; new as "new"; } }
