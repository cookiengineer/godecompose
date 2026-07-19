arch x86_64;
platform linux, darwin, windows, freebsd;

pattern time_Since {
    name: "time.Since";
    library: "go-stdlib";
    description: "Returns time since a given time";

    instr match {
        MOVQ t, AX
        CALL time_Since
    }

    gen {
        time.Since($t);
    }

    bind {
        t as "t";
    }
}
