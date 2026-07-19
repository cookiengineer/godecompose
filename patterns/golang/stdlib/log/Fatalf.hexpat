arch x86_64;
platform linux, darwin, windows, freebsd;

pattern log_Fatalf {
    name: "log.Fatalf";
    library: "go-stdlib";
    description: "Logs a formatted message and calls os.Exit(1)";

    instr match {
        CALL log_Fatalf
    }

    gen {
        log.Fatalf($format, $args...);
    }

    bind {
        format as "fmt";
        args as "args";
    }
}
