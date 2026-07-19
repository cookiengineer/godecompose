arch x86_64;
platform linux, darwin, windows, freebsd;

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
