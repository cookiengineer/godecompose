arch x86_64;
platform linux, darwin, windows, freebsd;

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
