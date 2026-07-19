arch x86_64;
platform linux, darwin, windows, freebsd;

pattern fmt_Fprintf {
    name: "fmt.Fprintf";
    library: "go-stdlib";
    description: "Formats and writes to an io.Writer";

    instr match {
        MOVQ w, AX
        MOVQ fmt_str, BX
        LEAQ args, CX
        CALL fmt_Fprintf
    }

    gen {
        fmt.Fprintf($w, $format, $args...);
    }

    bind {
        w as "w";
        fmt_str as "format";
        args as "args";
    }
}
