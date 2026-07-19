arch x86_64;
platform linux, darwin, windows, freebsd;

pattern sync_Mutex_Unlock_call {
    name: "sync.Mutex.Unlock";
    library: "go-stdlib";
    description: "sync.Mutex.Unlock call in ABIInternal";

    instr match {
        CALL sync_Mutex_Unlock
    }

    gen {
        $mu.Unlock();
    }

    bind {
        mu as "mu";
    }
}
