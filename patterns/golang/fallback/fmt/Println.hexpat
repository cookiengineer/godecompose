arch x86_64;
platform linux, darwin, windows, freebsd;

pattern fmt_Println_string {
    name: "fmt.Println";
    library: "go-stdlib";
    description: "fmt.Println with a string argument — recovers to fmt.Println(str)";

    instr match {
        CALL fmt_Fprintln
    }

    gen {
        fmt.Println($str);
    }

    bind {
        str as "str";
    }
}
