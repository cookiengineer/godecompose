arch x86_64;
platform linux, darwin, windows, freebsd;

pattern bytes_LastIndex {
    name: "bytes.LastIndex"; library: "go-stdlib"; instr match { CALL bytes_LastIndex }
    gen { bytes.LastIndex($s, $sep); } bind { s as "s"; sep as "sep"; }
}
