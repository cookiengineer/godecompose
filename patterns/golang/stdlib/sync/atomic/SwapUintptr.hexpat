arch x86_64;
platform linux, darwin, windows, freebsd;

pattern atomic_SwapUintptr { name: "sync/atomic.SwapUintptr"; library: "go-stdlib"; instr match { CALL sync_atomic_SwapUintptr } gen { atomic.SwapUintptr($addr, $new); } bind { addr as "addr"; new as "new"; } }
