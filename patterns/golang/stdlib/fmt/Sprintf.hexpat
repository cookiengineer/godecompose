arch x86_64;
platform linux, darwin, windows, freebsd;

pattern fmt_Sprintf {
    name: "fmt.Sprintf";
    library: "go-stdlib";
    description: "Formats and returns a string";

    instr match {
        MOVQ fmt_str, AX
        LEAQ args, BX
        CALL fmt_Sprintf
    }

    gen {
        fmt.Sprintf($format, $args...);
    }

    bind {
        fmt_str as "format";
        args as "args";
    }
}
