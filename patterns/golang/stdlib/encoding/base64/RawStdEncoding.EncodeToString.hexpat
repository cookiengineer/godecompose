arch x86_64;
platform linux, darwin, windows, freebsd;

pattern base64_RawStdEncoding_EncodeToString {
    name: "base64.RawStdEncoding.EncodeToString"; library: "go-stdlib";
    instr match { CALL encoding_base64___Encoding__EncodeToString }
    gen { base64.RawStdEncoding.EncodeToString($src); } bind { src as "src"; }
}
