arch x86_64;
platform linux, darwin, windows, freebsd;

pattern url_PathUnescape {
    name: "url.PathUnescape"; library: "go-stdlib";
    instr match { CALL net_url_PathUnescape }
    gen { url.PathUnescape($s); } bind { s as "s"; }
}
