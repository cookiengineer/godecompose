arch x86_64;
platform linux, darwin, windows, freebsd;

pattern atomic_StoreUint32 { name: "sync/atomic.StoreUint32"; library: "go-stdlib"; instr match { CALL sync_atomic_StoreUint32 } gen { atomic.StoreUint32($addr, $val); } bind { addr as "addr"; val as "val"; } }
