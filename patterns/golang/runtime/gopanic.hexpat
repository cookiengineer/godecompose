arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_gopanic {
    name: "runtime.gopanic";
    library: "go-runtime";
    description: "Initiates a panic";

    instr match {
        MOVQ eface, AX
        CALL runtime_gopanic
    }

    gen {
        panic($val);
    }

    bind {
        eface as "val";
    }
}
