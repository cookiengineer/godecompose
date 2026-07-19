arch x86_64;
platform linux, darwin, windows, freebsd;

pattern net_LookupHost {
    name: "net.LookupHost"; library: "go-stdlib";
    instr match { CALL net_LookupHost }
    gen { net.LookupHost($host); } bind { host as "host"; }
}
