arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_gcStart {
    name: "runtime.gcStart";
    library: "go-runtime";
    description: "Starts a GC cycle";

    instr match {
        MOVQ trigger, AX
        CALL runtime_gcStart
    }

    gen {
        gcStart($trigger);
    }

    bind {
        trigger as "trigger";
    }
}
