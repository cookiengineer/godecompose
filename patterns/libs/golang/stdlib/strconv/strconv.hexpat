arch x86_64;
platform linux, darwin, windows, freebsd;

pattern strconv_Itoa {
    name: "strconv.Itoa";
    library: "go-stdlib";
    description: "Converts an integer to a string";

    instr match {
        CALL strconv_Itoa
    }

    gen {
        strconv.Itoa($n);
    }

    bind {
        n as "n";
    }
}

pattern strconv_Atoi {
    name: "strconv.Atoi";
    library: "go-stdlib";
    description: "Converts a string to an integer";

    instr match {
        CALL strconv_Atoi
    }

    gen {
        strconv.Atoi($s);
    }

    bind {
        s as "s";
    }
}

pattern strconv_FormatInt {
    name: "strconv.FormatInt";
    library: "go-stdlib";
    description: "Formats an int64 as a string in given base";

    instr match {
        CALL strconv_FormatInt
    }

    gen {
        strconv.FormatInt($n, $base);
    }

    bind {
        n as "n";
        base as "base";
    }
}

pattern strconv_ParseInt {
    name: "strconv.ParseInt";
    library: "go-stdlib";
    description: "Parses a string to int64";

    instr match {
        CALL strconv_ParseInt
    }

    gen {
        strconv.ParseInt($s, $base, $bitSize);
    }

    bind {
        s as "s";
        base as "base";
        bitSize as "bitSize";
    }
}

pattern strconv_FormatFloat {
    name: "strconv.FormatFloat";
    library: "go-stdlib";
    description: "Formats a float as a string";

    instr match {
        CALL strconv_FormatFloat
    }

    gen {
        strconv.FormatFloat($f, $fmt, $prec, $bitSize);
    }

    bind {
        f as "f";
        fmt as "fmt";
        prec as "prec";
        bitSize as "bitSize";
    }
}

pattern strconv_ParseFloat {
    name: "strconv.ParseFloat";
    library: "go-stdlib";
    description: "Parses a string to float64";

    instr match {
        CALL strconv_ParseFloat
    }

    gen {
        strconv.ParseFloat($s, $bitSize);
    }

    bind {
        s as "s";
        bitSize as "bitSize";
    }
}

pattern strconv_Quote {
    name: "strconv.Quote";
    library: "go-stdlib";
    description: "Returns a Go string literal";

    instr match {
        CALL strconv_Quote
    }

    gen {
        strconv.Quote($s);
    }

    bind {
        s as "s";
    }
}

pattern strconv_Unquote {
    name: "strconv.Unquote";
    library: "go-stdlib";
    description: "Unquotes a Go string literal";

    instr match {
        CALL strconv_Unquote
    }

    gen {
        strconv.Unquote($s);
    }

    bind {
        s as "s";
    }
}
