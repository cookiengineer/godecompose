arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_makemap {
    name: "runtime.makemap";
    library: "go-runtime";
    description: "Creates a new map";

    instr match {
        MOVQ typ, AX
        MOVQ hint, BX
        MOVQ hmap_buf, CX
        CALL runtime_makemap
    }

    gen {
        make(map[$key_type]$val_type, $hint);
    }

    bind {
        typ as "type";
        hint as "size";
    }
}
