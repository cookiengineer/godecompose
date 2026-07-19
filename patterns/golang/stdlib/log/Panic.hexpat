arch x86_64;
platform linux, darwin, windows, freebsd;

pattern log_Panic {
    name: "log.Panic";
    library: "go-stdlib";
    description: "Logs a message and panics";

    instr match {
        CALL log_Panic
    }

    gen {
        log.Panic($args...);
    }

    bind {
        args as "args";
    }
}
