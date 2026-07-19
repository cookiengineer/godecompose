arch x86_64;
platform linux, darwin, windows, freebsd;

pattern fmt_Sprintf_format {
    name: "fmt.Sprintf";
    library: "go-stdlib";
    description: "fmt.Sprintf with format string and args — recovers to fmt.Sprintf(fmt, args...)";

    instr match {
        CALL fmt_Sprintf
    }

    gen {
        s := fmt.Sprintf($format, $args...);
    }

    bind {
        format as "fmt";
        args as "args";
    }
}
