arch x86_64;
platform linux, darwin, windows, freebsd;

pattern hex_NewDecoder {
    name: "hex.NewDecoder"; library: "go-stdlib";
    instr match { CALL encoding_hex_NewDecoder }
    gen { hex.NewDecoder($r); } bind { r as "r"; }
}
