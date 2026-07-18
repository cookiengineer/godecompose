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

pattern sync_WaitGroup_Add {
    name: "sync.WaitGroup.Add";
    library: "go-stdlib";
    description: "Adds delta to WaitGroup counter";

    instr match {
        MOVQ wg, AX
        MOVQ delta, BX
        CALL sync_WaitGroup_Add
    }

    gen {
        $wg.Add($delta);
    }

    bind {
        wg as "wg";
        delta as "delta";
    }
}

pattern sync_WaitGroup_Done {
    name: "sync.WaitGroup.Done";
    library: "go-stdlib";
    description: "Decrements WaitGroup counter";

    instr match {
        MOVQ wg, AX
        MOVQ minus_one, BX
        CALL sync_WaitGroup_Add
    }

    gen {
        $wg.Done();
    }

    bind {
        wg as "wg";
    }
}

pattern sync_WaitGroup_Wait {
    name: "sync.WaitGroup.Wait";
    library: "go-stdlib";
    description: "Waits for WaitGroup counter to reach zero";

    instr match {
        MOVQ wg, AX
        CALL sync_WaitGroup_Wait
    }

    gen {
        $wg.Wait();
    }

    bind {
        wg as "wg";
    }
}

pattern sync_Once_Do {
    name: "sync.Once.Do";
    library: "go-stdlib";
    description: "Executes function exactly once";

    instr match {
        MOVQ once, AX
        MOVQ fn, BX
        CALL sync_Once_Do
    }

    gen {
        $once.Do($fn);
    }

    bind {
        once as "once";
        fn as "fn";
    }
}

pattern sync_Pool_Get {
    name: "sync.Pool.Get";
    library: "go-stdlib";
    description: "Gets an item from a sync.Pool";

    instr match {
        MOVQ pool, AX
        CALL sync_Pool_Get
    }

    gen {
        $pool.Get();
    }

    bind {
        pool as "pool";
    }
}

pattern sync_Pool_Put {
    name: "sync.Pool.Put";
    library: "go-stdlib";
    description: "Puts an item back into a sync.Pool";

    instr match {
        MOVQ pool, AX
        MOVQ val, BX
        CALL sync_Pool_Put
    }

    gen {
        $pool.Put($val);
    }

    bind {
        pool as "pool";
        val as "val";
    }
}
