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

pattern fmt_Println {
    name: "fmt.Println";
    library: "go-stdlib";
    description: "Prints arguments with newline";

    instr match {
        LEAQ args, AX
        CALL fmt_Println
    }

    gen {
        fmt.Println($args...);
    }

    bind {
        args as "args";
    }
}

pattern fmt_Fprint {
    name: "fmt.Fprint";
    library: "go-stdlib";
    description: "Writes arguments to an io.Writer";

    instr match {
        MOVQ w, AX
        LEAQ args, BX
        CALL fmt_Fprint
    }

    gen {
        fmt.Fprint($w, $args...);
    }

    bind {
        w as "w";
        args as "args";
    }
}

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
