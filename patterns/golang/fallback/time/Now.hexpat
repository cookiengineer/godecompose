arch x86_64;
platform linux, darwin, windows, freebsd;

pattern time_Now_call {
    name: "time.Now";
    library: "go-stdlib";
    description: "time.Now call";

    instr match {
        CALL time_Now
    }

    gen {
        now := time.Now();
    }
}
