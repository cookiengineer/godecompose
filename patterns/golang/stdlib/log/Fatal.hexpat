arch x86_64;
platform linux, darwin, windows, freebsd;

pattern log_Fatal {
    name: "log.Fatal";
    library: "go-stdlib";
    description: "Logs a message and calls os.Exit(1)";

    instr match {
        CALL log_Fatal
    }

    gen {
        log.Fatal($args...);
    }

    bind {
        args as "args";
    }
}
