arch x86_64;
platform linux, darwin, windows, freebsd;

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
