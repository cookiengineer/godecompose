arch x86_64;
platform linux, darwin, windows, freebsd;

pattern os_File_Write_call {
    name: "os.File.Write";
    library: "go-stdlib";
    description: "os.File.Write method call";

    instr match {
        CALL os_File_Write
    }

    gen {
        os.File.Write($buf);
    }

    bind {
        buf as "buf";
    }
}
