arch x86_64;
platform linux, darwin, windows, freebsd;

pattern bytes_Buffer_Bytes {
    name: "bytes.Buffer.Bytes"; library: "go-stdlib"; instr match { CALL bytes___Buffer__Bytes }
    gen { $b.Bytes(); } bind { b as "b"; }
}
