arch x86_64;
platform linux, darwin, windows, freebsd;

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
