arch x86_64;
platform linux, darwin, windows, freebsd;

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
