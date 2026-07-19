arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_makemap_small {
    name: "runtime.makemap_small";
    library: "go-runtime";
    description: "Creates a small map (no hint)";

    instr match {
        MOVQ typ, AX
        CALL runtime_makemap_small
    }

    gen {
        make(map[$key_type]$val_type);
    }

    bind {
        typ as "type";
    }
}
