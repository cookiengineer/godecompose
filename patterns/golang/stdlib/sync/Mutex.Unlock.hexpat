arch x86_64;
platform linux, darwin, windows, freebsd;

pattern sync_Mutex_Unlock {
    name: "sync.Mutex.Unlock";
    library: "go-stdlib";
    description: "Unlocks a sync.Mutex";

    instr match {
        MOVQ m, AX
        CALL sync_Mutex_Unlock
    }

    gen {
        $m.Unlock();
    }

    bind {
        m as "mu";
    }
}
