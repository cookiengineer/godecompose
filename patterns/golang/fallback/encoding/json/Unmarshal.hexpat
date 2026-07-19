arch x86_64;
platform linux, darwin, windows, freebsd;

pattern json_Unmarshal_call {
    name: "encoding/json.Unmarshal";
    library: "go-stdlib";
    description: "json.Unmarshal call";

    instr match {
        CALL json_Unmarshal
    }

    gen {
        err := json.Unmarshal($data, $v);
    }

    bind {
        data as "data";
        v as "v";
    }
}
