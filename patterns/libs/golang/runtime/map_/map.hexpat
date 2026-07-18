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

pattern runtime_mapdelete {
    name: "runtime.mapdelete";
    library: "go-runtime";
    description: "Deletes a key from a map";

    instr match {
        MOVQ typ, AX
        MOVQ hmap, BX
        MOVQ key, CX
        CALL runtime_mapdelete
    }

    gen {
        delete($hmap, $key);
    }

    bind {
        typ as "type";
        hmap as "m";
        key as "key";
    }
}

pattern runtime_mapiterinit {
    name: "runtime.mapiterinit";
    library: "go-runtime";
    description: "Initializes a map iterator";

    instr match {
        MOVQ typ, AX
        MOVQ hmap, BX
        LEAQ it_ptr, CX
        CALL runtime_mapiterinit
    }

    gen {
        for k := range $hmap { ...
    }

    bind {
        typ as "type";
        hmap as "m";
        it_ptr as "it";
    }
}

pattern runtime_mapiternext {
    name: "runtime.mapiternext";
    library: "go-runtime";
    description: "Advances a map iterator";

    instr match {
        LEAQ it_ptr, AX
        CALL runtime_mapiternext
    }

    gen {
        // next iteration
    }

    bind {
        it_ptr as "it";
    }
}

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
