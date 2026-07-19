arch x86_64;
platform linux, darwin, windows, freebsd;

pattern atomic_StoreUint64 { name: "sync/atomic.StoreUint64"; library: "go-stdlib"; instr match { CALL sync_atomic_StoreUint64 } gen { atomic.StoreUint64($addr, $val); } bind { addr as "addr"; val as "val"; } }
