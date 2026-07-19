arch x86_64;
platform linux, darwin, windows, freebsd;

pattern net_ResolveUDPAddr {
    name: "net.ResolveUDPAddr"; library: "go-stdlib";
    instr match { CALL net_ResolveUDPAddr }
    gen { net.ResolveUDPAddr($network, $address); } bind { network as "network"; address as "address"; }
}
