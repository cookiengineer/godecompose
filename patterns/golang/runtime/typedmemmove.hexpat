arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_typedmemmove {
    name: "runtime.typedmemmove";
    library: "go-runtime";
    description: "Copies typed memory with write barrier";

    instr match {
        MOVQ typ, AX
        MOVQ dst, BX
        MOVQ src, CX
        CALL runtime_typedmemmove
    }

    gen {
        typedmemmove($dst, $src);
    }

    bind {
        typ as "type";
        dst as "dst";
        src as "src";
    }
}
