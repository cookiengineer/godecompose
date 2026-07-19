arch x86_64;
platform linux, darwin, windows, freebsd;

pattern net_InterfaceByName {
    name: "net.InterfaceByName"; library: "go-stdlib";
    instr match { CALL net_InterfaceByName }
    gen { net.InterfaceByName($name); } bind { name as "name"; }
}
