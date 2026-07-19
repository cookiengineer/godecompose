arch x86_64;
platform linux, darwin, windows, freebsd;

pattern http_ListenAndServeTLS {
    name: "http.ListenAndServeTLS"; library: "go-stdlib";
    instr match { CALL net_http_ListenAndServeTLS }
    gen { http.ListenAndServeTLS($addr, $certFile, $keyFile, $handler); } bind { addr as "addr"; certFile as "certFile"; keyFile as "keyFile"; handler as "handler"; }
}
