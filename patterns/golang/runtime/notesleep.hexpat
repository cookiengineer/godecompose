arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_notesleep {
    name: "runtime.notesleep";
    library: "go-runtime";
    description: "Sleeps on a runtime note";

    instr match {
        MOVQ n, AX
        CALL runtime_notesleep
    }

    gen {
        notesleep($n);
    }

    bind {
        n as "n";
    }
}
