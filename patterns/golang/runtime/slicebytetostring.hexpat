arch x86_64;
platform linux, darwin, windows, freebsd;

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
