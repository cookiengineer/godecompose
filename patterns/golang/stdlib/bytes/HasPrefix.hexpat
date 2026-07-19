arch x86_64;
platform linux, darwin, windows, freebsd;

pattern bytes_HasPrefix {
    name: "bytes.HasPrefix"; library: "go-stdlib"; instr match { CALL bytes_HasPrefix }
    gen { bytes.HasPrefix($s, $prefix); } bind { s as "s"; prefix as "prefix"; }
}
