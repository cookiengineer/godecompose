arch x86_64;
platform linux, darwin, windows, freebsd;

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
