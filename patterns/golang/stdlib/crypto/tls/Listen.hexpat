arch x86_64;
platform linux, darwin, windows, freebsd;

pattern tls_Listen {
    name: "tls.Listen"; library: "go-stdlib";
    instr match { CALL crypto_tls_Listen }
    gen { tls.Listen($network, $laddr, $config); } bind { network as "network"; laddr as "laddr"; config as "config"; }
}
