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
