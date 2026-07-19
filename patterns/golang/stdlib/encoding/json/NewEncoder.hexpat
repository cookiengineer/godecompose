arch x86_64;
platform linux, darwin, windows, freebsd;

pattern json_NewEncoder {
    name: "encoding/json.NewEncoder";
    library: "go-stdlib";
    description: "Creates a new JSON encoder";

    instr match {
        MOVQ w, AX
        CALL json_NewEncoder
    }

    gen {
        json.NewEncoder($w);
    }

    bind {
        w as "w";
    }
}
