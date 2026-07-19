arch x86_64;
platform linux, darwin, windows, freebsd;

pattern go_runtime_memmove {
    name: "runtime.memmove";
    library: "go-runtime";
    version: ">=1.0";
    description: "Go runtime memmove — copies n bytes from src to dst";

    instr match {
        MOVQ src, SI
        MOVQ dst, DI
        MOVQ len, CX
        REP; MOVSQ
        RET
    }

    gen {
        memmove($dst, $src, $len);
    }

    bind {
        src as "src";
        dst as "dst";
        len as "count";
    }
}
