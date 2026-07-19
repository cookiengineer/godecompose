arch x86_64;
platform linux, darwin, windows, freebsd;

pattern flag_FlagSet_Int {
    name: "flag.FlagSet.Int";
    library: "go-stdlib";
    description: "Defines an int flag";

    instr match {
        CALL flag___FlagSet__Int
    }

    gen {
        fs.Int($name, $value, $usage);
    }

    bind {
        name as "name";
        value as "value";
        usage as "usage";
    }
}
