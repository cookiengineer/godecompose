arch x86_64;
platform linux, darwin, windows, freebsd;

pattern io_Copy {
    name: "io.Copy";
    library: "go-stdlib";
    description: "Copies from reader to writer";

    instr match {
        MOVQ dst, AX
        MOVQ src, BX
        CALL io_Copy
    }

    gen {
        io.Copy($dst, $src);
    }

    bind {
        dst as "dst";
        src as "src";
    }
}
