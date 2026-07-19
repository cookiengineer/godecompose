arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_goexit {
    name: "runtime.goexit";
    library: "go-runtime";
    description: "Terminates the current goroutine";

    instr match {
        MOVL zero, AX
        CALL runtime_goexit1
    }

    gen {
        runtime_goexit();
    }
}
