arch x86_64;
platform linux, darwin, windows, freebsd;

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
