arch x86_64;
platform linux, darwin, windows, freebsd;

pattern net_DialTimeout {
    name: "net.DialTimeout"; library: "go-stdlib";
    instr match { CALL net_DialTimeout }
    gen { net.DialTimeout($network, $address, $timeout); } bind { network as "network"; address as "address"; timeout as "timeout"; }
}
