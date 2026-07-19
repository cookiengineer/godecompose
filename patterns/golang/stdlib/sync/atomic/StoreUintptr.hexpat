arch x86_64;
platform linux, darwin, windows, freebsd;

pattern atomic_StoreUintptr { name: "sync/atomic.StoreUintptr"; library: "go-stdlib"; instr match { CALL sync_atomic_StoreUintptr } gen { atomic.StoreUintptr($addr, $val); } bind { addr as "addr"; val as "val"; } }
