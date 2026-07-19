arch x86_64;
platform linux, darwin, windows, freebsd;

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
