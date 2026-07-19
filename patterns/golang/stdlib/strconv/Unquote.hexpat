arch x86_64;
platform linux, darwin, windows, freebsd;

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
