arch x86_64;
platform linux, darwin, windows, freebsd;

pattern http_NewServeMux {
    name: "http.NewServeMux"; library: "go-stdlib";
    instr match { CALL net_http_NewServeMux }
    gen { http.NewServeMux(); }
}
