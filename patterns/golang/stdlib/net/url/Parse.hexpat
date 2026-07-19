arch x86_64;
platform linux, darwin, windows, freebsd;

pattern url_Parse {
    name: "url.Parse"; library: "go-stdlib";
    instr match { CALL net_url_Parse }
    gen { url.Parse($rawURL); } bind { rawURL as "rawURL"; }
}
