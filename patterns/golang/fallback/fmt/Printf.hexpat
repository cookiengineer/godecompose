arch x86_64;
platform linux, darwin, windows, freebsd;

pattern fmt_Printf_format {
    name: "fmt.Printf";
    library: "go-stdlib";
    description: "fmt.Printf with format string and args — recovers to fmt.Printf(fmt, args...)";

    instr match {
        CALL fmt_Fprintf
    }

    gen {
        fmt.Printf($format, $args...);
    }

    bind {
        format as "fmt";
        args as "args";
    }
}
