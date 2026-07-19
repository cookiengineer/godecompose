arch x86_64;
platform linux, darwin, windows, freebsd;

pattern flag_Parse {
    name: "flag.Parse";
    library: "go-stdlib";
    description: "Parses command-line flags";

    instr match {
        CALL flag_Parse
    }

    gen {
        flag.Parse();
    }
}
