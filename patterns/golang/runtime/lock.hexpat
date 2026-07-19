arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_lock {
    name: "runtime.lock";
    library: "go-runtime";
    description: "Acquires a runtime mutex";

    instr match {
        MOVQ lk, AX
        CALL runtime_lock
    }

    gen {
        runtime_lock($lk);
    }

    bind {
        lk as "lk";
    }
}
