arch x86_64;
platform linux, darwin, windows, freebsd;

pattern tls_DialWithDialer {
    name: "tls.DialWithDialer"; library: "go-stdlib";
    instr match { CALL crypto_tls_DialWithDialer }
    gen { tls.DialWithDialer($dialer, $network, $addr, $config); } bind { dialer as "dialer"; network as "network"; addr as "addr"; config as "config"; }
}
