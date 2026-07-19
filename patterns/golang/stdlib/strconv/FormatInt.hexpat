arch x86_64;
platform linux, darwin, windows, freebsd;

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
