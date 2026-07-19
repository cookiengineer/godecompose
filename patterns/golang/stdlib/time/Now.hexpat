arch x86_64;
platform linux, darwin, windows, freebsd;

pattern time_Now {
    name: "time.Now";
    library: "go-stdlib";
    description: "Returns current time";

    instr match {
        CALL time_Now
    }

    gen {
        time.Now();
    }
}
