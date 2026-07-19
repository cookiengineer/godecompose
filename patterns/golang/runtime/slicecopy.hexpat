arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_slicecopy {
    name: "runtime.slicecopy";
    library: "go-runtime";
    description: "Copies elements between slices";

    instr match {
        MOVQ to_ptr, AX
        MOVQ to_len, BX
        MOVQ fm_ptr, CX
        MOVQ fm_len, DI
        MOVQ width, SI
        CALL runtime_slicecopy
    }

    gen {
        copy($dst, $src);
    }

    bind {
        to_ptr as "dst";
        fm_ptr as "src";
        width as "elemSize";
    }
}
