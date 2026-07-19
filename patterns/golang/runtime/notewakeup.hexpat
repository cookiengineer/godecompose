arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_notewakeup {
    name: "runtime.notewakeup";
    library: "go-runtime";
    description: "Wakes up a runtime note";

    instr match {
        MOVQ n, AX
        CALL runtime_notewakeup
    }

    gen {
        notewakeup($n);
    }

    bind {
        n as "n";
    }
}
