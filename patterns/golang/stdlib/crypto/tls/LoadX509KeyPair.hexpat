arch x86_64;
platform linux, darwin, windows, freebsd;

pattern tls_LoadX509KeyPair {
    name: "tls.LoadX509KeyPair"; library: "go-stdlib";
    instr match { CALL crypto_tls_LoadX509KeyPair }
    gen { tls.LoadX509KeyPair($certFile, $keyFile); } bind { certFile as "certFile"; keyFile as "keyFile"; }
}
