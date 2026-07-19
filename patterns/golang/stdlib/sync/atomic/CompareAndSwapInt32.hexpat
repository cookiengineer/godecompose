arch x86_64;
platform linux, darwin, windows, freebsd;

pattern atomic_CompareAndSwapInt32 { name: "sync/atomic.CompareAndSwapInt32"; library: "go-stdlib"; instr match { CALL sync_atomic_CompareAndSwapInt32 } gen { atomic.CompareAndSwapInt32($addr, $old, $new); } bind { addr as "addr"; old as "old"; new as "new"; } }
