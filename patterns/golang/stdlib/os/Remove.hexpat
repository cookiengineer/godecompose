arch x86_64;
platform linux, darwin, windows, freebsd;

pattern os_Remove {
    name: "os.Remove";
    library: "go-stdlib";
    description: "Removes a file";

    instr match {
        MOVQ name, AX
        CALL os_Remove
    }

    gen {
        os.Remove($name);
    }

    bind {
        name as "name";
    }
}
