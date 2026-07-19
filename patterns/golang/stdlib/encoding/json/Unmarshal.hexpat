arch x86_64;
platform linux, darwin, windows, freebsd;

pattern json_Unmarshal {
    name: "encoding/json.Unmarshal";
    library: "go-stdlib";
    description: "Unmarshals JSON into a value";

    instr match {
        MOVQ data, AX
        MOVQ val, BX
        CALL json_Unmarshal
    }

    gen {
        json.Unmarshal($data, $val);
    }

    bind {
        data as "data";
        val as "v";
    }
}
