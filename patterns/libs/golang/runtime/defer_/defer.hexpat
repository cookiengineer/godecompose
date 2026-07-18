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

pattern runtime_deferprocStack {
    name: "runtime.deferprocStack";
    library: "go-runtime";
    description: "Sets up a stacked deferred function call";

    instr match {
        LEAQ d, AX
        CALL runtime_deferprocStack
    }

    gen {
        defer $fn();
    }

    bind {
        d as "d";
    }
}

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

pattern runtime_fatal {
    name: "runtime.fatal";
    library: "go-runtime";
    description: "Fatal error, exits the program";

    instr match {
        MOVQ msg, AX
        CALL runtime_fatal
    }

    gen {
        panic($msg);
    }

    bind {
        msg as "msg";
    }
}

pattern runtime_throw {
    name: "runtime.throw";
    library: "go-runtime";
    description: "Throws a fatal runtime error";

    instr match {
        MOVQ msg, AX
        CALL runtime_throw
    }

    gen {
        panic($msg);
    }

    bind {
        msg as "msg";
    }
}
