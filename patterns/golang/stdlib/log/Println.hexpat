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
