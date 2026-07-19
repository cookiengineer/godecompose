arch x86_64;
platform linux, darwin, windows, freebsd;

pattern flag_FlagSet_String {
    name: "flag.FlagSet.String";
    library: "go-stdlib";
    description: "Defines a string flag";

    instr match {
        CALL flag___FlagSet__String
    }

    gen {
        fs.String($name, $value, $usage);
    }

    bind {
        name as "name";
        value as "value";
        usage as "usage";
    }
}
