arch x86_64;
platform linux, darwin, windows, freebsd;

pattern fmt_Errorf {
    name: "fmt.Errorf";
    library: "go-stdlib";
    description: "Formats an error string";

    instr match {
        MOVQ fmt_str, AX
        LEAQ args, BX
        CALL fmt_Errorf
    }

    gen {
        fmt.Errorf($format, $args...);
    }

    bind {
        fmt_str as "format";
        args as "args";
    }
}
