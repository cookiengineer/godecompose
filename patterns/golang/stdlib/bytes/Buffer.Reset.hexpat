arch x86_64;
platform linux, darwin, windows, freebsd;

pattern bytes_Buffer_Reset {
    name: "bytes.Buffer.Reset"; library: "go-stdlib"; instr match { CALL bytes___Buffer__Reset }
    gen { $b.Reset(); } bind { b as "b"; }
}
