arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_convT64 {
    name: "runtime.convT64";
    library: "go-runtime";
    description: "Converts uint64 to interface";

    instr match {
        MOVQ typ, AX
        MOVQ val, BX
        CALL runtime_convT64
    }

    gen {
        // convT64($val) -> interface
    }

    bind {
        typ as "type";
        val as "val";
    }
}
