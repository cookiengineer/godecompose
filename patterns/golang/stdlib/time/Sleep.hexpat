arch x86_64;
platform linux, darwin, windows, freebsd;

pattern time_Sleep {
    name: "time.Sleep";
    library: "go-stdlib";
    description: "Sleeps for a duration";

    instr match {
        MOVQ d, AX
        CALL time_Sleep
    }

    gen {
        time.Sleep($d);
    }

    bind {
        d as "d";
    }
}
