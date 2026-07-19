arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_mapaccess1 {
    name: "runtime.mapaccess1";
    library: "go-runtime";
    description: "Accesses a map entry, returns *value (does not return zero)";

    instr match {
        MOVQ typ, AX
        MOVQ hmap, BX
        MOVQ key, CX
        CALL runtime_mapaccess1
    }

    gen {
        $hmap[$key];
    }

    bind {
        typ as "type";
        hmap as "m";
        key as "key";
    }
}
