arch x86_64;
platform linux, darwin, windows, freebsd;

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
