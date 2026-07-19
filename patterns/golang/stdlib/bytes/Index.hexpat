arch x86_64;
platform linux, darwin, windows, freebsd;

pattern bytes_Index {
    name: "bytes.Index"; library: "go-stdlib"; instr match { CALL bytes_Index }
    gen { bytes.Index($s, $sep); } bind { s as "s"; sep as "sep"; }
}
