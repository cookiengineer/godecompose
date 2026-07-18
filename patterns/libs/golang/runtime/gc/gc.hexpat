arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_gcWriteBarrier {
    name: "runtime.gcWriteBarrier";
    library: "go-runtime";
    description: "GC write barrier check - flags write barrier buffer flush";

    instr match {
        CMPQ writeBarrier_flag, zero
        JNE @barrier_flush
        RET
    @barrier_flush:
    }

    gen {
        // GC write barrier check -> wbBufFlush
    }

    bind {
        writeBarrier_flag as "flag";
    }
}

pattern runtime_wbBufFlush {
    name: "runtime.wbBufFlush";
    library: "go-runtime";
    description: "Flushes the write barrier buffer";

    instr match {
        MOVQ runtime_wbBuf, AX
        CMPQ AX, BX
        JNE @flush
        RET
    @flush:
        MOVQ g_ptr, R14
        CALL runtime_wbBufFlush
    }
}

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

pattern runtime_gcDrain {
    name: "runtime.gcDrain";
    library: "go-runtime";
    description: "Drains GC work buffers";

    instr match {
        MOVQ gcw, AX
        CALL runtime_gcDrain
    }

    gen {
        gcDrain();
    }
}

pattern runtime_scanobject {
    name: "runtime.scanobject";
    library: "go-runtime";
    description: "Scans an object for GC pointers";

    instr match {
        MOVQ obj, AX
        MOVQ gcw, BX
        CALL runtime_scanobject
    }

    gen {
        scanobject($obj, $gcw);
    }

    bind {
        obj as "obj";
        gcw as "gcw";
    }
}

pattern runtime_greyobject {
    name: "runtime.greyobject";
    library: "go-runtime";
    description: "Marks an object grey during GC";

    instr match {
        MOVQ obj, AX
        MOVQ base, BX
        MOVQ off, CX
        MOVQ gcw, DI
        CALL runtime_greyobject
    }

    gen {
        greyobject($obj, $base, $off, $gcw);
    }

    bind {
        obj as "obj";
        base as "base";
        off as "off";
        gcw as "gcw";
    }
}

pattern runtime_gcStart {
    name: "runtime.gcStart";
    library: "go-runtime";
    description: "Starts a GC cycle";

    instr match {
        MOVQ trigger, AX
        CALL runtime_gcStart
    }

    gen {
        gcStart($trigger);
    }

    bind {
        trigger as "trigger";
    }
}
