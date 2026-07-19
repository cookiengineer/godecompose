arch x86_64;
platform linux, darwin, windows, freebsd;

pattern url_ParseRequestURI {
    name: "url.ParseRequestURI"; library: "go-stdlib";
    instr match { CALL net_url_ParseRequestURI }
    gen { url.ParseRequestURI($rawURL); } bind { rawURL as "rawURL"; }
}
