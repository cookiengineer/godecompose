arch x86_64;
platform linux, darwin, windows, freebsd;

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
