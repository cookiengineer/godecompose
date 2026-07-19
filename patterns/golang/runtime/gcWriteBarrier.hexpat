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
