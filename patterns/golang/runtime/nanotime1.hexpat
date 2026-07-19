arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_nanotime {
    name: "runtime.nanotime1";
    library: "go-runtime";
    description: "Gets current monotonic time in nanoseconds";

    instr match {
        CALL runtime_nanotime1
    }

    gen {
        runtime_nanotime();
    }
}
