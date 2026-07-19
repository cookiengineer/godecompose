arch x86_64;
platform linux, darwin, windows, freebsd;

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
