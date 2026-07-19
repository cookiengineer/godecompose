arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_deferproc {
    name: "runtime.deferproc";
    library: "go-runtime";
    description: "Sets up a deferred function call";

    instr match {
        MOVQ framesize, AX
        MOVQ fn_ptr, BX
        CALL runtime_deferproc
    }

    gen {
        defer $fn();
    }

    bind {
        framesize as "framesize";
        fn_ptr as "fn";
    }
}
