arch x86_64;
platform linux, darwin, windows, freebsd;

pattern tls_NewListener {
    name: "tls.NewListener"; library: "go-stdlib";
    instr match { CALL crypto_tls_NewListener }
    gen { tls.NewListener($inner, $config); } bind { inner as "inner"; config as "config"; }
}
