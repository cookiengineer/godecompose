arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_memequal {
    name: "runtime.memequal";
    library: "go-runtime";
    description: "Compares n bytes for equality";

    instr match {
        MOVQ a, SI
        MOVQ b, DI
        MOVQ len, CX
        SHRQ size3, CX
        REPE; CMPSQ
        JNE @notequal
        MOVQ len, CX
        ANDQ mask7, CX
        REPE; CMPSB
        JNE @notequal
        MOVQ one, AX
        RET
    @notequal:
        XORL AX, AX
        RET
    }

    gen {
        memequal($a, $b, $len) != 0;
    }

    bind {
        a as "a";
        b as "b";
        len as "n";
    }
}
