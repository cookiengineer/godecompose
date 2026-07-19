arch x86_64;
platform linux, darwin, windows, freebsd;

pattern atomic_StorePointer { name: "sync/atomic.StorePointer"; library: "go-stdlib"; instr match { CALL sync_atomic_StorePointer } gen { atomic.StorePointer($addr, $val); } bind { addr as "addr"; val as "val"; } }
