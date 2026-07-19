arch x86_64;
platform linux, darwin, windows, freebsd;

pattern http_Head {
    name: "http.Head"; library: "go-stdlib";
    instr match { CALL net_http_Head }
    gen { http.Head($url); } bind { url as "url"; }
}
