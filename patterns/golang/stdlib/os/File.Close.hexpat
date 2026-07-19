arch x86_64;
platform linux, darwin, windows, freebsd;

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
