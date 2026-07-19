arch x86_64;
platform linux, darwin, windows, freebsd;

pattern bytes_Join {
    name: "bytes.Join"; library: "go-stdlib"; instr match { CALL bytes_Join }
    gen { bytes.Join($s, $sep); } bind { s as "s"; sep as "sep"; }
}
