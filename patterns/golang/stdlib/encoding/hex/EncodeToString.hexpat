arch x86_64;
platform linux, darwin, windows, freebsd;

pattern hex_EncodeToString {
    name: "hex.EncodeToString"; library: "go-stdlib";
    instr match { CALL encoding_hex_EncodeToString }
    gen { hex.EncodeToString($src); } bind { src as "src"; }
}
