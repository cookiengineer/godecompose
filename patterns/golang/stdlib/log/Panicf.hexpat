arch x86_64;
platform linux, darwin, windows, freebsd;

pattern log_Panicf {
    name: "log.Panicf";
    library: "go-stdlib";
    description: "Logs a formatted message and panics";

    instr match {
        CALL log_Panicf
    }

    gen {
        log.Panicf($format, $args...);
    }

    bind {
        format as "fmt";
        args as "args";
    }
}
