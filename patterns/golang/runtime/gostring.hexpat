arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_gostring {
    name: "runtime.gostring";
    library: "go-runtime";
    description: "Creates a Go string from C string";

    instr match {
        MOVQ ptr, AX
        CALL runtime_gostring
    }

    gen {
        gostring($ptr);
    }

    bind {
        ptr as "ptr";
    }
}
