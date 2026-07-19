arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_typedmemclr {
    name: "runtime.typedmemclr";
    library: "go-runtime";
    description: "Clears typed memory with write barrier";

    instr match {
        MOVQ typ, AX
        MOVQ ptr, BX
        CALL runtime_typedmemclr
    }

    gen {
        typedmemclr($ptr);
    }

    bind {
        typ as "type";
        ptr as "ptr";
    }
}
