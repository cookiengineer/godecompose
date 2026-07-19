arch x86_64;
platform linux, darwin, windows, freebsd;

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
