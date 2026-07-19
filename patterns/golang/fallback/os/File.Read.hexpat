arch x86_64;
platform linux, darwin, windows, freebsd;

pattern os_File_Read_call {
    name: "os.File.Read";
    library: "go-stdlib";
    description: "os.File.Read method call";

    instr match {
        CALL os_File_Read
    }

    gen {
        os.File.Read($buf);
    }

    bind {
        buf as "buf";
    }
}
