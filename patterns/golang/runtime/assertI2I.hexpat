arch x86_64;
platform linux, darwin, windows, freebsd;

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
