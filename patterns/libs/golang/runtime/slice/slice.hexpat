arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_makeslice {
    name: "runtime.makeslice";
    library: "go-runtime";
    description: "Creates a new slice";

    instr match {
        MOVQ typ, AX
        MOVQ len, BX
        MOVQ cap, CX
        CALL runtime_makeslice
    }

    gen {
        make([]$typ, $len, $cap);
    }

    bind {
        typ as "type";
        len as "len";
        cap as "cap";
    }
}

pattern runtime_growslice {
    name: "runtime.growslice";
    library: "go-runtime";
    description: "Grows a slice capacity";

    instr match {
        MOVQ typ, AX
        MOVQ old_ptr, BX
        MOVQ new_cap, CX
        CALL runtime_growslice
    }

    gen {
        append($slice, ...);
    }

    bind {
        typ as "type";
        old_ptr as "slice";
        new_cap as "cap";
    }
}

pattern runtime_slicebytetostring {
    name: "runtime.slicebytetostring";
    library: "go-runtime";
    description: "Converts a byte slice to a string";

    instr match {
        MOVQ buf_ptr, AX
        MOVQ buf_len, BX
        CALL runtime_slicebytetostring
    }

    gen {
        string($buf);
    }

    bind {
        buf_ptr as "buf";
        buf_len as "len";
    }
}

pattern runtime_stringtoslicebyte {
    name: "runtime.stringtoslicebyte";
    library: "go-runtime";
    description: "Converts a string to a byte slice";

    instr match {
        MOVQ str_ptr, AX
        MOVQ str_len, BX
        CALL runtime_stringtoslicebyte
    }

    gen {
        []byte($str);
    }

    bind {
        str_ptr as "str";
        str_len as "len";
    }
}

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
