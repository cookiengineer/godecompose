arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_mapaccess2 {
    name: "runtime.mapaccess2";
    library: "go-runtime";
    description: "Accesses a map entry, returns *value and bool present";

    instr match {
        MOVQ typ, AX
        MOVQ hmap, BX
        MOVQ key, CX
        CALL runtime_mapaccess2
    }

    gen {
        _, ok := $hmap[$key];
    }

    bind {
        typ as "type";
        hmap as "m";
        key as "key";
    }
}
