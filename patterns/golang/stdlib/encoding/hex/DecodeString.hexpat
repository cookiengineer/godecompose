arch x86_64;
platform linux, darwin, windows, freebsd;

pattern hex_DecodeString {
    name: "hex.DecodeString"; library: "go-stdlib";
    instr match { CALL encoding_hex_DecodeString }
    gen { hex.DecodeString($s); } bind { s as "s"; }
}
