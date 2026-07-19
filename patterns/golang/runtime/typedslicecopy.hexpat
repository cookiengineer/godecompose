arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_typedslicecopy {
    name: "runtime.typedslicecopy";
    library: "go-runtime";
    description: "Copies elements between slices with type info";

    instr match {
        MOVQ typ, AX
        MOVQ dst_ptr, BX
        MOVQ dst_len, CX
        MOVQ src_ptr, DI
        MOVQ src_len, SI
        CALL runtime_typedslicecopy
    }

    gen {
        copy($dst, $src);
    }

    bind {
        dst_ptr as "dst";
        src_ptr as "src";
    }
}
