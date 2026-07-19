arch x86_64;
platform linux, darwin, windows, freebsd;

pattern http_PostForm {
    name: "http.PostForm"; library: "go-stdlib";
    instr match { CALL net_http_PostForm }
    gen { http.PostForm($url, $data); } bind { url as "url"; data as "data"; }
}
