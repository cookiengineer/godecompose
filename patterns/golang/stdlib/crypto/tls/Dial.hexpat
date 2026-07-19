arch x86_64;
platform linux, darwin, windows, freebsd;

pattern tls_Dial {
    name: "tls.Dial"; library: "go-stdlib";
    instr match { CALL crypto_tls_Dial }
    gen { tls.Dial($network, $addr, $config); } bind { network as "network"; addr as "addr"; config as "config"; }
}
