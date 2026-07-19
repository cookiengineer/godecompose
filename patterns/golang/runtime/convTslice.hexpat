arch x86_64;
platform linux, darwin, windows, freebsd;

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
