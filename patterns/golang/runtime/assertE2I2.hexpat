arch x86_64;
platform linux, darwin, windows, freebsd;

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
