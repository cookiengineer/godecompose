arch x86_64;
platform linux, darwin, windows, freebsd;

pattern time_After {
    name: "time.After";
    library: "go-stdlib";
    description: "Returns a channel that fires after a duration";

    instr match {
        MOVQ d, AX
        CALL time_After
    }

    gen {
        time.After($d);
    }

    bind {
        d as "d";
    }
}
