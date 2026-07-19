arch x86_64;
platform linux, darwin, windows, freebsd;

pattern os_Open {
    name: "os.Open";
    library: "go-stdlib";
    description: "Opens a file for reading";

    instr match {
        MOVQ name, AX
        CALL os_Open
    }

    gen {
        os.Open($name);
    }

    bind {
        name as "name";
    }
}
