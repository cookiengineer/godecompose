arch x86_64;
platform linux, darwin, windows, freebsd;

pattern os_Mkdir {
    name: "os.Mkdir";
    library: "go-stdlib";
    description: "Creates a directory";

    instr match {
        MOVQ name, AX
        MOVQ perm, BX
        CALL os_Mkdir
    }

    gen {
        os.Mkdir($name, $perm);
    }

    bind {
        name as "name";
        perm as "perm";
    }
}
