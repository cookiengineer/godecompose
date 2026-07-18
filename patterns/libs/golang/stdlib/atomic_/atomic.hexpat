arch x86_64;
platform linux, darwin, windows, freebsd;

pattern atomic_LoadInt32 { name: "sync/atomic.LoadInt32"; library: "go-stdlib"; instr match { CALL sync_atomic_LoadInt32 } gen { atomic.LoadInt32($addr); } bind { addr as "addr"; } }
pattern atomic_StoreInt32 { name: "sync/atomic.StoreInt32"; library: "go-stdlib"; instr match { CALL sync_atomic_StoreInt32 } gen { atomic.StoreInt32($addr, $val); } bind { addr as "addr"; val as "val"; } }
pattern atomic_AddInt32 { name: "sync/atomic.AddInt32"; library: "go-stdlib"; instr match { CALL sync_atomic_AddInt32 } gen { atomic.AddInt32($addr, $delta); } bind { addr as "addr"; delta as "delta"; } }
pattern atomic_SwapInt32 { name: "sync/atomic.SwapInt32"; library: "go-stdlib"; instr match { CALL sync_atomic_SwapInt32 } gen { atomic.SwapInt32($addr, $new); } bind { addr as "addr"; new as "new"; } }
pattern atomic_CompareAndSwapInt32 { name: "sync/atomic.CompareAndSwapInt32"; library: "go-stdlib"; instr match { CALL sync_atomic_CompareAndSwapInt32 } gen { atomic.CompareAndSwapInt32($addr, $old, $new); } bind { addr as "addr"; old as "old"; new as "new"; } }

pattern atomic_LoadInt64 { name: "sync/atomic.LoadInt64"; library: "go-stdlib"; instr match { CALL sync_atomic_LoadInt64 } gen { atomic.LoadInt64($addr); } bind { addr as "addr"; } }
pattern atomic_StoreInt64 { name: "sync/atomic.StoreInt64"; library: "go-stdlib"; instr match { CALL sync_atomic_StoreInt64 } gen { atomic.StoreInt64($addr, $val); } bind { addr as "addr"; val as "val"; } }
pattern atomic_AddInt64 { name: "sync/atomic.AddInt64"; library: "go-stdlib"; instr match { CALL sync_atomic_AddInt64 } gen { atomic.AddInt64($addr, $delta); } bind { addr as "addr"; delta as "delta"; } }
pattern atomic_SwapInt64 { name: "sync/atomic.SwapInt64"; library: "go-stdlib"; instr match { CALL sync_atomic_SwapInt64 } gen { atomic.SwapInt64($addr, $new); } bind { addr as "addr"; new as "new"; } }
pattern atomic_CompareAndSwapInt64 { name: "sync/atomic.CompareAndSwapInt64"; library: "go-stdlib"; instr match { CALL sync_atomic_CompareAndSwapInt64 } gen { atomic.CompareAndSwapInt64($addr, $old, $new); } bind { addr as "addr"; old as "old"; new as "new"; } }

pattern atomic_LoadUint32 { name: "sync/atomic.LoadUint32"; library: "go-stdlib"; instr match { CALL sync_atomic_LoadUint32 } gen { atomic.LoadUint32($addr); } bind { addr as "addr"; } }
pattern atomic_StoreUint32 { name: "sync/atomic.StoreUint32"; library: "go-stdlib"; instr match { CALL sync_atomic_StoreUint32 } gen { atomic.StoreUint32($addr, $val); } bind { addr as "addr"; val as "val"; } }
pattern atomic_AddUint32 { name: "sync/atomic.AddUint32"; library: "go-stdlib"; instr match { CALL sync_atomic_AddUint32 } gen { atomic.AddUint32($addr, $delta); } bind { addr as "addr"; delta as "delta"; } }
pattern atomic_SwapUint32 { name: "sync/atomic.SwapUint32"; library: "go-stdlib"; instr match { CALL sync_atomic_SwapUint32 } gen { atomic.SwapUint32($addr, $new); } bind { addr as "addr"; new as "new"; } }
pattern atomic_CompareAndSwapUint32 { name: "sync/atomic.CompareAndSwapUint32"; library: "go-stdlib"; instr match { CALL sync_atomic_CompareAndSwapUint32 } gen { atomic.CompareAndSwapUint32($addr, $old, $new); } bind { addr as "addr"; old as "old"; new as "new"; } }

pattern atomic_LoadUint64 { name: "sync/atomic.LoadUint64"; library: "go-stdlib"; instr match { CALL sync_atomic_LoadUint64 } gen { atomic.LoadUint64($addr); } bind { addr as "addr"; } }
pattern atomic_StoreUint64 { name: "sync/atomic.StoreUint64"; library: "go-stdlib"; instr match { CALL sync_atomic_StoreUint64 } gen { atomic.StoreUint64($addr, $val); } bind { addr as "addr"; val as "val"; } }
pattern atomic_AddUint64 { name: "sync/atomic.AddUint64"; library: "go-stdlib"; instr match { CALL sync_atomic_AddUint64 } gen { atomic.AddUint64($addr, $delta); } bind { addr as "addr"; delta as "delta"; } }
pattern atomic_SwapUint64 { name: "sync/atomic.SwapUint64"; library: "go-stdlib"; instr match { CALL sync_atomic_SwapUint64 } gen { atomic.SwapUint64($addr, $new); } bind { addr as "addr"; new as "new"; } }
pattern atomic_CompareAndSwapUint64 { name: "sync/atomic.CompareAndSwapUint64"; library: "go-stdlib"; instr match { CALL sync_atomic_CompareAndSwapUint64 } gen { atomic.CompareAndSwapUint64($addr, $old, $new); } bind { addr as "addr"; old as "old"; new as "new"; } }

pattern atomic_LoadUintptr { name: "sync/atomic.LoadUintptr"; library: "go-stdlib"; instr match { CALL sync_atomic_LoadUintptr } gen { atomic.LoadUintptr($addr); } bind { addr as "addr"; } }
pattern atomic_StoreUintptr { name: "sync/atomic.StoreUintptr"; library: "go-stdlib"; instr match { CALL sync_atomic_StoreUintptr } gen { atomic.StoreUintptr($addr, $val); } bind { addr as "addr"; val as "val"; } }
pattern atomic_AddUintptr { name: "sync/atomic.AddUintptr"; library: "go-stdlib"; instr match { CALL sync_atomic_AddUintptr } gen { atomic.AddUintptr($addr, $delta); } bind { addr as "addr"; delta as "delta"; } }
pattern atomic_SwapUintptr { name: "sync/atomic.SwapUintptr"; library: "go-stdlib"; instr match { CALL sync_atomic_SwapUintptr } gen { atomic.SwapUintptr($addr, $new); } bind { addr as "addr"; new as "new"; } }
pattern atomic_CompareAndSwapUintptr { name: "sync/atomic.CompareAndSwapUintptr"; library: "go-stdlib"; instr match { CALL sync_atomic_CompareAndSwapUintptr } gen { atomic.CompareAndSwapUintptr($addr, $old, $new); } bind { addr as "addr"; old as "old"; new as "new"; } }

pattern atomic_LoadPointer { name: "sync/atomic.LoadPointer"; library: "go-stdlib"; instr match { CALL sync_atomic_LoadPointer } gen { atomic.LoadPointer($addr); } bind { addr as "addr"; } }
pattern atomic_StorePointer { name: "sync/atomic.StorePointer"; library: "go-stdlib"; instr match { CALL sync_atomic_StorePointer } gen { atomic.StorePointer($addr, $val); } bind { addr as "addr"; val as "val"; } }
pattern atomic_SwapPointer { name: "sync/atomic.SwapPointer"; library: "go-stdlib"; instr match { CALL sync_atomic_SwapPointer } gen { atomic.SwapPointer($addr, $new); } bind { addr as "addr"; new as "new"; } }
pattern atomic_CompareAndSwapPointer { name: "sync/atomic.CompareAndSwapPointer"; library: "go-stdlib"; instr match { CALL sync_atomic_CompareAndSwapPointer } gen { atomic.CompareAndSwapPointer($addr, $old, $new); } bind { addr as "addr"; old as "old"; new as "new"; } }
