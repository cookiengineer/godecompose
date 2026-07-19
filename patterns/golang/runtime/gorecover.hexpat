arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_gorecover {
    name: "runtime.gorecover";
    library: "go-runtime";
    description: "Recovers from a panic";

    instr match {
        LEAQ ret, AX
        CALL runtime_gorecover
    }

    gen {
        recover();
    }

    bind {
        ret as "ret";
    }
}
