arch x86_64;
platform linux, darwin, windows, freebsd;

pattern net_LookupAddr {
    name: "net.LookupAddr"; library: "go-stdlib";
    instr match { CALL net_LookupAddr }
    gen { net.LookupAddr($addr); } bind { addr as "addr"; }
}
