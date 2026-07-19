arch x86_64;
platform linux, darwin, windows, freebsd;

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
