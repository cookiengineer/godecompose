arch x86_64;
platform linux, darwin, windows, freebsd;

pattern json_NewDecoder {
    name: "encoding/json.NewDecoder";
    library: "go-stdlib";
    description: "Creates a new JSON decoder";

    instr match {
        MOVQ r, AX
        CALL json_NewDecoder
    }

    gen {
        json.NewDecoder($r);
    }

    bind {
        r as "r";
    }
}
