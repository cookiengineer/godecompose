arch x86_64;
platform linux, darwin, windows, freebsd;

pattern os_Stat {
    name: "os.Stat";
    library: "go-stdlib";
    description: "Returns file info";

    instr match {
        MOVQ name, AX
        CALL os_Stat
    }

    gen {
        os.Stat($name);
    }

    bind {
        name as "name";
    }
}
