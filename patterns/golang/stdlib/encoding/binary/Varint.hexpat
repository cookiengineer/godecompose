arch x86_64;
platform linux, darwin, windows, freebsd;

pattern binary_Varint {
    name: "binary.Varint"; library: "go-stdlib";
    instr match { CALL encoding_binary_Varint }
    gen { binary.Varint($buf); } bind { buf as "buf"; }
}
