arch x86_64;
platform linux, darwin, windows, freebsd;

pattern atomic_StoreInt64 { name: "sync/atomic.StoreInt64"; library: "go-stdlib"; instr match { CALL sync_atomic_StoreInt64 } gen { atomic.StoreInt64($addr, $val); } bind { addr as "addr"; val as "val"; } }
