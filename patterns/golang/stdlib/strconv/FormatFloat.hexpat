arch x86_64;
platform linux, darwin, windows, freebsd;

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
