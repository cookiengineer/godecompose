arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_convTstring {
    name: "runtime.convTstring";
    library: "go-runtime";
    description: "Converts string to interface";

    instr match {
        MOVQ typ, AX
        MOVQ str_ptr, BX
        MOVQ str_len, CX
        CALL runtime_convTstring
    }

    gen {
        // convTstring($str) -> interface
    }

    bind {
        typ as "type";
        str_ptr as "str";
    }
}
