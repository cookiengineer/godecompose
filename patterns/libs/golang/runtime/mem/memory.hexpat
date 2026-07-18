arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_memmove {
    name: "runtime.memmove";
    library: "go-runtime";
    description: "Copies n bytes from src to dst, handles overlap";

    instr match {
        CMPQ src, dst
        JBE @forward
        STD
        MOVQ src, SI
        MOVQ len, CX
        MOVQ dst, DI
        ADDQ len, SI
        ADDQ len, DI
        SUBQ one, SI
        SUBQ one, DI
        REP; MOVSB
        CLD
        RET
    @forward:
        MOVQ src, SI
        MOVQ dst, DI
        MOVQ len, CX
        SHRQ size3, CX
        REP; MOVSQ
        MOVQ len, CX
        ANDQ mask7, CX
        REP; MOVSB
        RET
    }

    gen {
        memmove($dst, $src, $len);
    }

    bind {
        src as "src";
        dst as "dst";
        len as "n";
    }
}

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

pattern runtime_memequal {
    name: "runtime.memequal";
    library: "go-runtime";
    description: "Compares n bytes for equality";

    instr match {
        MOVQ a, SI
        MOVQ b, DI
        MOVQ len, CX
        SHRQ size3, CX
        REPE; CMPSQ
        JNE @notequal
        MOVQ len, CX
        ANDQ mask7, CX
        REPE; CMPSB
        JNE @notequal
        MOVQ one, AX
        RET
    @notequal:
        XORL AX, AX
        RET
    }

    gen {
        memequal($a, $b, $len) != 0;
    }

    bind {
        a as "a";
        b as "b";
        len as "n";
    }
}

pattern runtime_memhash {
    name: "runtime.memhash";
    library: "go-runtime";
    description: "Hashes n bytes of memory";

    instr match {
        MOVQ ptr, AX
        MOVQ len, BX
        MOVQ seed, CX
        CALL aeshashbody
    }

    gen {
        memhash($ptr, $seed, $len);
    }

    bind {
        ptr as "ptr";
        len as "size";
        seed as "seed";
    }
}

pattern runtime_newobject {
    name: "runtime.newobject";
    library: "go-runtime";
    description: "Allocates and zeroes a new object of the given type";

    instr match {
        MOVQ typ, AX
        CALL runtime_mallocgc
    }

    gen {
        new(*$typ)();
    }

    bind {
        typ as "type";
    }
}

pattern runtime_newarray {
    name: "runtime.newarray";
    library: "go-runtime";
    description: "Allocates an array of n elements of the given type";

    instr match {
        MOVQ typ, AX
        MOVQ count, BX
        CALL runtime_mallocgc
    }

    gen {
        make([]$typ, $count);
    }

    bind {
        typ as "type";
        count as "n";
    }
}

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
