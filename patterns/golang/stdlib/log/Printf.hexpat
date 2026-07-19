arch x86_64;
platform linux, darwin, windows, freebsd;

pattern log_Printf {
    name: "log.Printf";
    library: "go-stdlib";
    description: "Logs a formatted message";

    instr match {
        CALL log_Printf
    }

    gen {
        log.Printf($format, $args...);
    }

    bind {
        format as "fmt";
        args as "args";
    }
}
