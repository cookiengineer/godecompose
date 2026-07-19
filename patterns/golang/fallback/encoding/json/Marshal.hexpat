arch x86_64;
platform linux, darwin, windows, freebsd;

pattern json_Marshal_call {
    name: "encoding/json.Marshal";
    library: "go-stdlib";
    description: "json.Marshal call";

    instr match {
        CALL json_Marshal
    }

    gen {
        data, err := json.Marshal($v);
    }

    bind {
        v as "v";
    }
}
