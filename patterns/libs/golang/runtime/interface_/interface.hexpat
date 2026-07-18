arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_convT64 {
    name: "runtime.convT64";
    library: "go-runtime";
    description: "Converts uint64 to interface";

    instr match {
        MOVQ typ, AX
        MOVQ val, BX
        CALL runtime_convT64
    }

    gen {
        // convT64($val) -> interface
    }

    bind {
        typ as "type";
        val as "val";
    }
}

pattern runtime_convTstring {
    name: "runtime.convTstring";
    library: "go-runtime";
    description: "Converts string to interface";

    instr match {
        MOVQ typ, AX
        MOVQ str_ptr, BX
        MOVQ str_len, CX
        CALL runtime_convTstring
    }

    gen {
        // convTstring($str) -> interface
    }

    bind {
        typ as "type";
        str_ptr as "str";
    }
}

pattern runtime_convTslice {
    name: "runtime.convTslice";
    library: "go-runtime";
    description: "Converts slice to interface";

    instr match {
        MOVQ typ, AX
        MOVQ ptr, BX
        MOVQ len, CX
        MOVQ cap, DI
        CALL runtime_convTslice
    }

    gen {
        // convTslice($slice) -> interface
    }

    bind {
        typ as "type";
        ptr as "slice";
    }
}

pattern runtime_convT2E {
    name: "runtime.convT2E";
    library: "go-runtime";
    description: "Converts a type to an empty interface (eface)";

    instr match {
        MOVQ typ, AX
        MOVQ elem, BX
        CALL runtime_convT2E
    }

    gen {
        // convT2E($elem)
    }

    bind {
        typ as "type";
        elem as "val";
    }
}

pattern runtime_convT2I {
    name: "runtime.convT2I";
    library: "go-runtime";
    description: "Converts a type to a non-empty interface (iface)";

    instr match {
        MOVQ tab, AX
        MOVQ elem, BX
        CALL runtime_convT2I
    }

    gen {
        // convT2I($elem)
    }

    bind {
        tab as "itab";
        elem as "val";
    }
}

pattern runtime_assertE2I {
    name: "runtime.assertE2I";
    library: "go-runtime";
    description: "Asserts empty interface to non-empty interface";

    instr match {
        MOVQ tab, AX
        MOVQ eface, BX
        CALL runtime_assertE2I
    }

    gen {
        // $eface.(Itab)
    }

    bind {
        tab as "itab";
        eface as "val";
    }
}

pattern runtime_assertI2I {
    name: "runtime.assertI2I";
    library: "go-runtime";
    description: "Asserts non-empty interface to another interface";

    instr match {
        MOVQ tab, AX
        MOVQ iface, BX
        CALL runtime_assertI2I
    }

    gen {
        // $iface.(Itab)
    }

    bind {
        tab as "itab";
        iface as "val";
    }
}

pattern runtime_assertE2I2 {
    name: "runtime.assertE2I2";
    library: "go-runtime";
    description: "Asserts with ok flag";

    instr match {
        MOVQ tab, AX
        MOVQ eface, BX
        CALL runtime_assertE2I2
    }

    gen {
        // val, ok := $eface.(Itab)
    }

    bind {
        tab as "itab";
        eface as "val";
    }
}

pattern runtime_ifaceeq {
    name: "runtime.ifaceeq";
    library: "go-runtime";
    description: "Compares two non-empty interfaces";

    instr match {
        MOVQ iface1, AX
        MOVQ iface2, BX
        CALL runtime_ifaceeq
    }

    gen {
        // ifaceeq($a, $b)
    }

    bind {
        iface1 as "a";
        iface2 as "b";
    }
}

pattern runtime_efaceeq {
    name: "runtime.efaceeq";
    library: "go-runtime";
    description: "Compares two empty interfaces";

    instr match {
        MOVQ eface1, AX
        MOVQ eface2, BX
        CALL runtime_efaceeq
    }

    gen {
        // efaceeq($a, $b)
    }

    bind {
        eface1 as "a";
        eface2 as "b";
    }
}

pattern runtime_nilinterhash {
    name: "runtime.nilinterhash";
    library: "go-runtime";
    description: "Hashes a value for nil interface map key";

    instr match {
        MOVQ eface, AX
        MOVQ seed, BX
        CALL runtime_nilinterhash
    }

    gen {
        // nilinterhash($val, $seed)
    }

    bind {
        eface as "val";
        seed as "seed";
    }
}

pattern runtime_typehash {
    name: "runtime.typehash";
    library: "go-runtime";
    description: "Hashes a value with type info";

    instr match {
        MOVQ typ, AX
        MOVQ val, BX
        MOVQ seed, CX
        CALL runtime_typehash
    }

    gen {
        // typehash($val, $seed)
    }

    bind {
        typ as "type";
        val as "val";
        seed as "seed";
    }
}
