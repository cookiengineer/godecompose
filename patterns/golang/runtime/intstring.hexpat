arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_intstring {
    name: "runtime.intstring";
    library: "go-runtime";
    description: "Converts an integer to a string";

    instr match {
        LEAQ buf, AX
        MOVQ val, BX
        CALL runtime_intstring
    }

    gen {
        string(rune($val));
    }

    bind {
        val as "r";
    }
}
