arch x86_64;
platform linux, darwin, windows, freebsd;

pattern atomic_StoreInt32 { name: "sync/atomic.StoreInt32"; library: "go-stdlib"; instr match { CALL sync_atomic_StoreInt32 } gen { atomic.StoreInt32($addr, $val); } bind { addr as "addr"; val as "val"; } }
