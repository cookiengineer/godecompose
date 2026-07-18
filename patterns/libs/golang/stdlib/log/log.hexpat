arch x86_64;
platform linux, darwin, windows, freebsd;

pattern log_Println {
    name: "log.Println";
    library: "go-stdlib";
    description: "Logs a message with newline";

    instr match {
        CALL log_Println
    }

    gen {
        log.Println($args...);
    }

    bind {
        args as "args";
    }
}

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
