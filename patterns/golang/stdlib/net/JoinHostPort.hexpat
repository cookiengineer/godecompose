arch x86_64;
platform linux, darwin, windows, freebsd;

pattern net_JoinHostPort {
    name: "net.JoinHostPort"; library: "go-stdlib";
    instr match { CALL net_JoinHostPort }
    gen { net.JoinHostPort($host, $port); } bind { host as "host"; port as "port"; }
}
