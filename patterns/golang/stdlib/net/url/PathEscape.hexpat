arch x86_64;
platform linux, darwin, windows, freebsd;

pattern url_PathEscape {
    name: "url.PathEscape"; library: "go-stdlib";
    instr match { CALL net_url_PathEscape }
    gen { url.PathEscape($s); } bind { s as "s"; }
}
