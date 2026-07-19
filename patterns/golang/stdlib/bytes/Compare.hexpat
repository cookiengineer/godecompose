arch x86_64;
platform linux, darwin, windows, freebsd;

pattern bytes_Compare {
    name: "bytes.Compare"; library: "go-stdlib"; instr match { CALL bytes_Compare }
    gen { bytes.Compare($a, $b); } bind { a as "a"; b as "b"; }
}
