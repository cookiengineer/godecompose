arch x86_64;
platform linux, darwin, windows, freebsd;

pattern io_WriteString {
    name: "io.WriteString";
    library: "go-stdlib";
    description: "Writes a string to a writer";

    instr match {
        MOVQ w, AX
        MOVQ s, BX
        CALL io_WriteString
    }

    gen {
        io.WriteString($w, $s);
    }

    bind {
        w as "w";
        s as "s";
    }
}
