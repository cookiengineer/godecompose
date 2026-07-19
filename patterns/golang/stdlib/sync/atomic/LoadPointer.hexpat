arch x86_64;
platform linux, darwin, windows, freebsd;

pattern atomic_LoadPointer { name: "sync/atomic.LoadPointer"; library: "go-stdlib"; instr match { CALL sync_atomic_LoadPointer } gen { atomic.LoadPointer($addr); } bind { addr as "addr"; } }
