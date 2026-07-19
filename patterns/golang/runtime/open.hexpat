arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_open {
    name: "runtime.open";
    library: "go-runtime";
    description: "Opens a file";

    instr match {
        MOVQ name, AX
        MOVQ flags, BX
        MOVQ mode, CX
        CALL runtime_open
    }

    gen {
        runtime_open($name, $flags, $mode);
    }

    bind {
        name as "name";
        flags as "flags";
        mode as "mode";
    }
}
