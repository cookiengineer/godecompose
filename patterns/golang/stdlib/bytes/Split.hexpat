arch x86_64;
platform linux, darwin, windows, freebsd;

pattern bytes_Split {
    name: "bytes.Split"; library: "go-stdlib"; instr match { CALL bytes_Split }
    gen { bytes.Split($s, $sep); } bind { s as "s"; sep as "sep"; }
}
