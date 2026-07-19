arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_mapassign {
    name: "runtime.mapassign";
    library: "go-runtime";
    description: "Assigns a key in a map, returns pointer to value slot";

    instr match {
        MOVQ typ, AX
        MOVQ hmap, BX
        MOVQ key, CX
        CALL runtime_mapassign
    }

    gen {
        $hmap[$key] = ...;
    }

    bind {
        typ as "type";
        hmap as "m";
        key as "key";
    }
}
