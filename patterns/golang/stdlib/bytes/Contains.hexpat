arch x86_64;
platform linux, darwin, windows, freebsd;

pattern bytes_Contains {
    name: "bytes.Contains"; library: "go-stdlib"; instr match { CALL bytes_Contains }
    gen { bytes.Contains($b, $subslice); } bind { b as "b"; subslice as "subslice"; }
}
