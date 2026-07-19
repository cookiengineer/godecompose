arch x86_64;
platform linux, darwin, windows, freebsd;

pattern binary_PutUvarint {
    name: "binary.PutUvarint"; library: "go-stdlib";
    instr match { CALL encoding_binary_PutUvarint }
    gen { binary.PutUvarint($buf, $x); } bind { buf as "buf"; x as "x"; }
}
