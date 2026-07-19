arch x86_64;
platform linux, darwin, windows, freebsd;

pattern bytes_NewBuffer {
    name: "bytes.NewBuffer"; library: "go-stdlib"; instr match { CALL bytes_NewBuffer }
    gen { bytes.NewBuffer($buf); } bind { buf as "buf"; }
}
