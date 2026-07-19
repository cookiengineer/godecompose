arch x86_64;
platform linux, darwin, windows, freebsd;

pattern bytes_Buffer_WriteString {
    name: "bytes.Buffer.WriteString"; library: "go-stdlib"; instr match { CALL bytes___Buffer__WriteString }
    gen { $b.WriteString($s); } bind { b as "b"; s as "s"; }
}
