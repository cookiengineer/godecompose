arch x86_64;
platform linux, darwin, windows, freebsd;

pattern atomic_CompareAndSwapPointer { name: "sync/atomic.CompareAndSwapPointer"; library: "go-stdlib"; instr match { CALL sync_atomic_CompareAndSwapPointer } gen { atomic.CompareAndSwapPointer($addr, $old, $new); } bind { addr as "addr"; old as "old"; new as "new"; } }
