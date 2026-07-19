arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_newproc {
    name: "runtime.newproc";
    library: "go-runtime";
    description: "Creates a new goroutine";

    instr match {
        MOVQ fn_ptr, AX
        CALL runtime_newproc
    }

    gen {
        go $fn_ptr;
    }

    bind {
        fn_ptr as "fn";
    }
}
