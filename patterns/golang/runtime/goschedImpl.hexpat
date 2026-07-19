arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_gosched {
    name: "runtime.goschedImpl";
    library: "go-runtime";
    description: "Yields the processor";

    instr match {
        CALL runtime_goschedImpl
    }

    gen {
        runtime.Gosched();
    }
}
