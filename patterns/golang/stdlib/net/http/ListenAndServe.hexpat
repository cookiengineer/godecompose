arch x86_64;
platform linux, darwin, windows, freebsd;

pattern http_ListenAndServe {
    name: "http.ListenAndServe"; library: "go-stdlib";
    instr match { CALL net_http_ListenAndServe }
    gen { http.ListenAndServe($addr, $handler); } bind { addr as "addr"; handler as "handler"; }
}
