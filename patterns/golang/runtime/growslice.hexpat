arch x86_64;
platform linux, darwin, windows, freebsd;

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
