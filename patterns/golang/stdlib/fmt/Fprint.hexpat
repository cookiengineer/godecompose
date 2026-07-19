arch x86_64;
platform linux, darwin, windows, freebsd;

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
