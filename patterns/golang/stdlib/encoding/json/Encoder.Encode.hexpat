arch x86_64;
platform linux, darwin, windows, freebsd;

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
