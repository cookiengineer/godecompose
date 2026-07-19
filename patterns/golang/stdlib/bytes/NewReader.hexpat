arch x86_64;
platform linux, darwin, windows, freebsd;

pattern bytes_NewReader {
    name: "bytes.NewReader"; library: "go-stdlib"; instr match { CALL bytes_NewReader }
    gen { bytes.NewReader($b); } bind { b as "b"; }
}
