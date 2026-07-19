arch x86_64;
platform linux, darwin, windows, freebsd;

pattern fmt_Println {
    name: "fmt.Println";
    library: "go-stdlib";
    description: "Prints arguments with newline";

    instr match {
        LEAQ args, AX
        CALL fmt_Println
    }

    gen {
        fmt.Println($args...);
    }

    bind {
        args as "args";
    }
}
