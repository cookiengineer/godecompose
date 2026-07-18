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

pattern json_Encoder_Encode {
    name: "encoding/json.Encoder.Encode";
    library: "go-stdlib";
    description: "Encodes a value as JSON";

    instr match {
        MOVQ enc, AX
        MOVQ val, BX
        CALL json_Encoder_Encode
    }

    gen {
        $enc.Encode($val);
    }

    bind {
        enc as "enc";
        val as "v";
    }
}

pattern json_Decoder_Decode {
    name: "encoding/json.Decoder.Decode";
    library: "go-stdlib";
    description: "Decodes JSON into a value";

    instr match {
        MOVQ dec, AX
        MOVQ val, BX
        CALL json_Decoder_Decode
    }

    gen {
        $dec.Decode($val);
    }

    bind {
        dec as "dec";
        val as "v";
    }
}
