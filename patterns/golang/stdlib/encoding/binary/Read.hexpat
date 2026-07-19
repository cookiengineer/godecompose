arch x86_64;
platform linux, darwin, windows, freebsd;

pattern binary_Read {
    name: "binary.Read"; library: "go-stdlib";
    instr match { CALL encoding_binary_Read }
    gen { binary.Read($r, $order, $data); } bind { r as "r"; order as "order"; data as "data"; }
}
