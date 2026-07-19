arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_deferreturn {
    name: "runtime.deferreturn";
    library: "go-runtime";
    description: "Runs deferred functions before returning";

    instr match {
        MOVQ framesize, AX
        CALL runtime_deferreturn
    }

    gen {
        // run deferred functions
    }

    bind {
        framesize as "framesize";
    }
}
