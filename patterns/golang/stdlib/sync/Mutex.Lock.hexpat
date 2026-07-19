arch x86_64;
platform linux, darwin, windows, freebsd;

pattern sync_Mutex_Lock {
    name: "sync.Mutex.Lock";
    library: "go-stdlib";
    description: "Locks a sync.Mutex";

    instr match {
        MOVQ m, AX
        CALL sync_Mutex_Lock
    }

    gen {
        $m.Lock();
    }

    bind {
        m as "mu";
    }
}
