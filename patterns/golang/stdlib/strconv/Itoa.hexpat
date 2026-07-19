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
