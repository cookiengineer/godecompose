arch x86_64;
platform linux, darwin, windows, freebsd;

pattern binary_Write {
    name: "binary.Write"; library: "go-stdlib";
    instr match { CALL encoding_binary_Write }
    gen { binary.Write($w, $order, $data); } bind { w as "w"; order as "order"; data as "data"; }
}
