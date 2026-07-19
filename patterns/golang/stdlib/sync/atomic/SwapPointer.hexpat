arch x86_64;
platform linux, darwin, windows, freebsd;

pattern atomic_SwapPointer { name: "sync/atomic.SwapPointer"; library: "go-stdlib"; instr match { CALL sync_atomic_SwapPointer } gen { atomic.SwapPointer($addr, $new); } bind { addr as "addr"; new as "new"; } }
