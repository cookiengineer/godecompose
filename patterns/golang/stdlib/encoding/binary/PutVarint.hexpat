arch x86_64;
platform linux, darwin, windows, freebsd;

pattern binary_PutVarint {
    name: "binary.PutVarint"; library: "go-stdlib";
    instr match { CALL encoding_binary_PutVarint }
    gen { binary.PutVarint($buf, $x); } bind { buf as "buf"; x as "x"; }
}
