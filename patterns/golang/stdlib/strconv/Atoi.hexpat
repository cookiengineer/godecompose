arch x86_64;
platform linux, darwin, windows, freebsd;

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
