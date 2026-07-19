arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_reflectcallmove {
    name: "runtime.reflectcallmove";
    library: "go-runtime";
    description: "Moves call arguments for reflection";

    instr match {
        MOVQ typ, AX
        MOVQ src, BX
        MOVQ dst, CX
        CALL runtime_reflectcallmove
    }

    gen {
        // reflectcallmove($typ, $src, $dst)
    }

    bind {
        typ as "typ";
        src as "src";
        dst as "dst";
    }
}
