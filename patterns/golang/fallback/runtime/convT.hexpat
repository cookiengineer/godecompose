arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_convT {
    name: "runtime.convT";
    library: "go-runtime";
    description: "Type conversion to interface";

    instr match {
        CALL runtime_convT64
        | CALL runtime_convTstring
        | CALL runtime_convTslice
        | CALL runtime_convT2E
        | CALL runtime_convT2I
    }

    gen {
        // convT: convert to interface
    }
}
