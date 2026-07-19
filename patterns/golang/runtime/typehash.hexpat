arch x86_64;
platform linux, darwin, windows, freebsd;

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
