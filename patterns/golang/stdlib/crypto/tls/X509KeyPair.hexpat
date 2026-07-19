arch x86_64;
platform linux, darwin, windows, freebsd;

pattern tls_X509KeyPair {
    name: "tls.X509KeyPair"; library: "go-stdlib";
    instr match { CALL crypto_tls_X509KeyPair }
    gen { tls.X509KeyPair($certPEMBlock, $keyPEMBlock); } bind { certPEMBlock as "certPEMBlock"; keyPEMBlock as "keyPEMBlock"; }
}
