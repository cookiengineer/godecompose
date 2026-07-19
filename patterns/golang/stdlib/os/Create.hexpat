arch x86_64;
platform linux, darwin, windows, freebsd;

pattern os_Create {
    name: "os.Create";
    library: "go-stdlib";
    description: "Creates or truncates a file";

    instr match {
        MOVQ name, AX
        CALL os_Create
    }

    gen {
        os.Create($name);
    }

    bind {
        name as "name";
    }
}
