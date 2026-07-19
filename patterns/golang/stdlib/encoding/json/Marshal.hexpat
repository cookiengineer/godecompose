arch x86_64;
platform linux, darwin, windows, freebsd;

pattern json_Marshal {
    name: "encoding/json.Marshal";
    library: "go-stdlib";
    description: "Marshals a value to JSON";

    instr match {
        MOVQ val, AX
        CALL json_Marshal
    }

    gen {
        json.Marshal($val);
    }

    bind {
        val as "v";
    }
}
