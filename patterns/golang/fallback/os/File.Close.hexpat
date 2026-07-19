arch x86_64;
platform linux, darwin, windows, freebsd;

pattern os_File_Close_call {
    name: "os.File.Close";
    library: "go-stdlib";
    description: "os.File.Close method call";

    instr match {
        CALL os_File_Close
    }

    gen {
        os.File.Close();
    }
}
