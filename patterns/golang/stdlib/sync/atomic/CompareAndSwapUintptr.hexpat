arch x86_64;
platform linux, darwin, windows, freebsd;

pattern atomic_CompareAndSwapUintptr { name: "sync/atomic.CompareAndSwapUintptr"; library: "go-stdlib"; instr match { CALL sync_atomic_CompareAndSwapUintptr } gen { atomic.CompareAndSwapUintptr($addr, $old, $new); } bind { addr as "addr"; old as "old"; new as "new"; } }
