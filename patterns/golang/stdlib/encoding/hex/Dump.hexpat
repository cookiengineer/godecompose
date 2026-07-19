arch x86_64;
platform linux, darwin, windows, freebsd;

pattern hex_Dump {
    name: "hex.Dump"; library: "go-stdlib";
    instr match { CALL encoding_hex_Dump }
    gen { hex.Dump($data); } bind { data as "data"; }
}
