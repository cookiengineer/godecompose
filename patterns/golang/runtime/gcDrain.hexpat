arch x86_64;
platform linux, darwin, windows, freebsd;

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
