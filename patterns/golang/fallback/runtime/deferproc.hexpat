arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_defer {
    name: "runtime.deferproc";
    library: "go-runtime";
    description: "Defer statement setup";

    instr match {
        CALL runtime_deferproc
        | CALL runtime_deferprocStack
    }

    gen {
        defer func();
    }
}
