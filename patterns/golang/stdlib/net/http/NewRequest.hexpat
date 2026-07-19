arch x86_64;
platform linux, darwin, windows, freebsd;

pattern http_NewRequest {
    name: "http.NewRequest"; library: "go-stdlib";
    instr match { CALL net_http_NewRequest }
    gen { http.NewRequest($method, $url, $body); } bind { method as "method"; url as "url"; body as "body"; }
}
