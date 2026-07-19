arch x86_64;
platform linux, darwin, windows, freebsd;

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
