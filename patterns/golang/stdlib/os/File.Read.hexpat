arch x86_64;
platform linux, darwin, windows, freebsd;

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
