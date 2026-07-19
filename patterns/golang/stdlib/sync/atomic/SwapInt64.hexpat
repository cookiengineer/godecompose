arch x86_64;
platform linux, darwin, windows, freebsd;

pattern atomic_SwapInt64 { name: "sync/atomic.SwapInt64"; library: "go-stdlib"; instr match { CALL sync_atomic_SwapInt64 } gen { atomic.SwapInt64($addr, $new); } bind { addr as "addr"; new as "new"; } }
