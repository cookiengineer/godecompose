arch x86_64;
platform linux, darwin, windows, freebsd;

pattern net_ResolveIPAddr {
    name: "net.ResolveIPAddr"; library: "go-stdlib";
    instr match { CALL net_ResolveIPAddr }
    gen { net.ResolveIPAddr($network, $address); } bind { network as "network"; address as "address"; }
}
