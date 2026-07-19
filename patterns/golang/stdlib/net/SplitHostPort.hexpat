arch x86_64;
platform linux, darwin, windows, freebsd;

pattern net_SplitHostPort {
    name: "net.SplitHostPort"; library: "go-stdlib";
    instr match { CALL net_SplitHostPort }
    gen { net.SplitHostPort($hostport); } bind { hostport as "hostport"; }
}
