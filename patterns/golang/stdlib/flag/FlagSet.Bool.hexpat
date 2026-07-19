arch x86_64;
platform linux, darwin, windows, freebsd;

pattern flag_FlagSet_Bool {
    name: "flag.FlagSet.Bool";
    library: "go-stdlib";
    description: "Defines a bool flag";

    instr match {
        CALL flag___FlagSet__Bool
    }

    gen {
        fs.Bool($name, $value, $usage);
    }

    bind {
        name as "name";
        value as "value";
        usage as "usage";
    }
}
