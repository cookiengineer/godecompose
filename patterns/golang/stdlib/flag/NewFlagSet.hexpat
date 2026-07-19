arch x86_64;
platform linux, darwin, windows, freebsd;

pattern flag_NewFlagSet {
    name: "flag.NewFlagSet";
    library: "go-stdlib";
    description: "Creates a new flag set";

    instr match {
        CALL flag_NewFlagSet
    }

    gen {
        fs := flag.NewFlagSet($name, $errorHandling);
    }

    bind {
        name as "name";
        errorHandling as "errorHandling";
    }
}
