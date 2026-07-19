arch x86_64;
platform linux, darwin, windows, freebsd;

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
