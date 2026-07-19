arch x86_64;
platform linux, darwin, windows, freebsd;

pattern sync_Mutex_Lock_call {
    name: "sync.Mutex.Lock";
    library: "go-stdlib";
    description: "sync.Mutex.Lock call in ABIInternal";

    instr match {
        CALL sync_Mutex_Lock
    }

    gen {
        $mu.Lock();
    }

    bind {
        mu as "mu";
    }
}
