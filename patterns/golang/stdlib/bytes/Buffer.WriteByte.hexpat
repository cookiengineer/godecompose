arch x86_64;
platform linux, darwin, windows, freebsd;

pattern bytes_Buffer_WriteByte {
    name: "bytes.Buffer.WriteByte"; library: "go-stdlib"; instr match { CALL bytes___Buffer__WriteByte }
    gen { $b.WriteByte($c); } bind { b as "b"; c as "c"; }
}
