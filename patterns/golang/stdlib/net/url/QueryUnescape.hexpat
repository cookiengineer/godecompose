arch x86_64;
platform linux, darwin, windows, freebsd;

pattern url_QueryUnescape {
    name: "url.QueryUnescape"; library: "go-stdlib";
    instr match { CALL net_url_QueryUnescape }
    gen { url.QueryUnescape($s); } bind { s as "s"; }
}
