arch x86_64;
platform linux, darwin, windows, freebsd;

pattern atomic_LoadInt64 { name: "sync/atomic.LoadInt64"; library: "go-stdlib"; instr match { CALL sync_atomic_LoadInt64 } gen { atomic.LoadInt64($addr); } bind { addr as "addr"; } }
