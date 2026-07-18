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

pattern flag_FlagSet_Parse {
    name: "flag.FlagSet.Parse";
    library: "go-stdlib";
    description: "Parses flag definitions from the argument list";

    instr match {
        CALL flag___FlagSet__Parse
    }

    gen {
        $fs.Parse($arguments);
    }

    bind {
        fs as "fs";
        arguments as "arguments";
    }
}
