arch x86_64;
platform linux, darwin, windows, freebsd;

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
