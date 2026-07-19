arch x86_64;
platform linux, darwin, windows, freebsd;

pattern atomic_CompareAndSwapInt64 { name: "sync/atomic.CompareAndSwapInt64"; library: "go-stdlib"; instr match { CALL sync_atomic_CompareAndSwapInt64 } gen { atomic.CompareAndSwapInt64($addr, $old, $new); } bind { addr as "addr"; old as "old"; new as "new"; } }
