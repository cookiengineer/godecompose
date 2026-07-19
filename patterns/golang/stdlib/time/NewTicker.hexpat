arch x86_64;
platform linux, darwin, windows, freebsd;

pattern time_Tick {
    name: "time.NewTicker";
    library: "go-stdlib";
    description: "Creates a new ticker";

    instr match {
        MOVQ d, AX
        CALL time_NewTicker
    }

    gen {
        time.NewTicker($d);
    }

    bind {
        d as "d";
    }
}
