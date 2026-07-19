arch x86_64;
platform linux, darwin, windows, freebsd;

pattern net_InterfaceAddrs {
    name: "net.InterfaceAddrs"; library: "go-stdlib";
    instr match { CALL net_InterfaceAddrs }
    gen { net.InterfaceAddrs(); }
}
