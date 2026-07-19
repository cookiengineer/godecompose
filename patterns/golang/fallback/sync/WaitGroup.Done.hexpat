arch x86_64;
platform linux, darwin, windows, freebsd;

pattern sync_WaitGroup_Done_call {
    name: "sync.WaitGroup.Done";
    library: "go-stdlib";
    description: "sync.WaitGroup.Done via Add(-1)";

    instr match {
        CALL sync_WaitGroup_Add
    }

    gen {
        $wg.Done();
    }

    bind {
        wg as "wg";
    }
}
