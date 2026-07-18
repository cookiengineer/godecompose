arch x86_64;
platform linux, darwin, windows, freebsd;

pattern os_File_Write {
    name: "os.File.Write";
    library: "go-stdlib";
    description: "Writes bytes to a file";

    instr match {
        MOVQ f, AX
        MOVQ buf, BX
        CALL os_File_Write
    }

    gen {
        $f.Write($buf);
    }

    bind {
        f as "f";
        buf as "buf";
    }
}

pattern os_File_Read {
    name: "os.File.Read";
    library: "go-stdlib";
    description: "Reads bytes from a file";

    instr match {
        MOVQ f, AX
        MOVQ buf, BX
        CALL os_File_Read
    }

    gen {
        $f.Read($buf);
    }

    bind {
        f as "f";
        buf as "buf";
    }
}

pattern os_File_Close {
    name: "os.File.Close";
    library: "go-stdlib";
    description: "Closes a file";

    instr match {
        MOVQ f, AX
        CALL os_File_Close
    }

    gen {
        $f.Close();
    }

    bind {
        f as "f";
    }
}

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

pattern os_Getenv {
    name: "os.Getenv";
    library: "go-stdlib";
    description: "Gets an environment variable";

    instr match {
        MOVQ key, AX
        CALL os_Getenv
    }

    gen {
        os.Getenv($key);
    }

    bind {
        key as "key";
    }
}
