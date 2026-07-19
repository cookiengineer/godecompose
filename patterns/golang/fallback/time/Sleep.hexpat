arch x86_64;
platform linux, darwin, windows, freebsd;

pattern time_Sleep_call {
    name: "time.Sleep";
    library: "go-stdlib";
    description: "time.Sleep call";

    instr match {
        CALL time_Sleep
    }

    gen {
        time.Sleep($d);
    }

    bind {
        d as "d";
    }
}
