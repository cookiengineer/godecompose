arch x86_64;
platform linux, darwin, windows, freebsd;

pattern bytes_Equal {
    name: "bytes.Equal"; library: "go-stdlib"; instr match { CALL bytes_Equal }
    gen { bytes.Equal($a, $b); } bind { a as "a"; b as "b"; }
}
