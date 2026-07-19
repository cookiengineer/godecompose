arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_memclrNoHeapPointers {
    name: "runtime.memclrNoHeapPointers";
    library: "go-runtime";
    description: "Clears n bytes to zero with no GC write barrier";

    instr match {
        XORQ AX, AX
        MOVQ ptr, DI
        MOVQ len, CX
        SHRQ size3, CX
        REP; STOSQ
        MOVQ len, CX
        ANDQ mask7, CX
        REP; STOSB
        RET
    }

    gen {
        memset($ptr, 0, $len);
    }

    bind {
        ptr as "ptr";
        len as "n";
    }
}
