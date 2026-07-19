arch x86_64;
platform linux, darwin, windows, freebsd;

pattern base64_StdEncoding_DecodeString {
    name: "base64.StdEncoding.DecodeString"; library: "go-stdlib";
    instr match { CALL encoding_base64___Encoding__DecodeString }
    gen { base64.StdEncoding.DecodeString($s); } bind { s as "s"; }
}
