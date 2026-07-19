arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_mallocgc {
    name: "runtime.mallocgc";
    library: "go-runtime";
    description: "Core allocation function";

    instr match {
        MOVQ size, AX
        MOVQ typ, BX
        MOVQ flags, CX
        CALL runtime_mallocgc
    }

    gen {
        mallocgc($size, $typ, $flags);
    }

    bind {
        size as "size";
        typ as "typ";
        flags as "flags";
    }
}
