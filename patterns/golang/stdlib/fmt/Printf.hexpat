arch x86_64;
platform linux, darwin, windows, freebsd;

pattern fmt_Printf {
    name: "fmt.Printf";
    library: "go-stdlib";
    description: "Formats and writes to stdout";

    instr match {
        MOVQ fmt_str, AX
        LEAQ args, BX
        CALL fmt_Printf
    }

    gen {
        fmt.Printf($format, $args...);
    }

    bind {
        fmt_str as "format";
        args as "args";
    }
}
