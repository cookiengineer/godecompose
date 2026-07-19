arch x86_64;
platform linux, darwin, windows, freebsd;

pattern bytes_Buffer_String {
    name: "bytes.Buffer.String"; library: "go-stdlib"; instr match { CALL bytes___Buffer__String }
    gen { $b.String(); } bind { b as "b"; }
}
