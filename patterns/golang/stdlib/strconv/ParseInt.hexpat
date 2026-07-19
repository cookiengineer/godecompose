arch x86_64;
platform linux, darwin, windows, freebsd;

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
