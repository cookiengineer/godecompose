arch x86_64;
platform linux, darwin, windows, freebsd;

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
