arch x86_64;
platform linux, darwin, windows, freebsd;

pattern io_ReadAll {
    name: "io.ReadAll";
    library: "go-stdlib";
    description: "Reads everything from a reader";

    instr match {
        MOVQ r, AX
        CALL io_ReadAll
    }

    gen {
        io.ReadAll($r);
    }

    bind {
        r as "r";
    }
}
