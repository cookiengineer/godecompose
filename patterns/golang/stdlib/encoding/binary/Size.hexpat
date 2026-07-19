arch x86_64;
platform linux, darwin, windows, freebsd;

pattern binary_Size {
    name: "binary.Size"; library: "go-stdlib";
    instr match { CALL encoding_binary_Size }
    gen { binary.Size($v); } bind { v as "v"; }
}
