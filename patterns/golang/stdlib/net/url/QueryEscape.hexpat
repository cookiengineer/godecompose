arch x86_64;
platform linux, darwin, windows, freebsd;

pattern url_QueryEscape {
    name: "url.QueryEscape"; library: "go-stdlib";
    instr match { CALL net_url_QueryEscape }
    gen { url.QueryEscape($s); } bind { s as "s"; }
}
