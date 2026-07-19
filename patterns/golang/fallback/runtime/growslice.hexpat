arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_growslice_call {
    name: "runtime.growslice";
    library: "go-runtime";
    description: "runtime.growslice call — append operation";

    instr match {
        CALL runtime_growslice
    }

    gen {
        append($s, ...);
    }

    bind {
        s as "s";
    }
}
