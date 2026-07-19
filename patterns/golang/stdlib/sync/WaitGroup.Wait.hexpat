arch x86_64;
platform linux, darwin, windows, freebsd;

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
