arch x86_64;
platform linux, darwin, windows, freebsd;

pattern net_ListenPacket {
    name: "net.ListenPacket"; library: "go-stdlib";
    instr match { CALL net_ListenPacket }
    gen { net.ListenPacket($network, $address); } bind { network as "network"; address as "address"; }
}
