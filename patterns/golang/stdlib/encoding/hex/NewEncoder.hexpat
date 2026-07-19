arch x86_64;
platform linux, darwin, windows, freebsd;

pattern hex_NewEncoder {
    name: "hex.NewEncoder"; library: "go-stdlib";
    instr match { CALL encoding_hex_NewEncoder }
    gen { hex.NewEncoder($w); } bind { w as "w"; }
}
