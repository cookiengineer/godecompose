arch x86_64;
platform linux, darwin, windows, freebsd;

pattern fmt_Errorf_format {
    name: "fmt.Errorf";
    library: "go-stdlib";
    description: "fmt.Errorf with format string and args — recovers to fmt.Errorf(fmt, args...)";

    instr match {
        CALL fmt_Errorf
    }

    gen {
        err := fmt.Errorf($format, $args...);
    }

    bind {
        format as "fmt";
        args as "args";
    }
}
