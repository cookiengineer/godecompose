arch x86_64;
platform linux, darwin, windows, freebsd;

pattern http_Request_WithContext {
    name: "http.Request.WithContext"; library: "go-stdlib";
    instr match { CALL net_http___Request__WithContext }
    gen { $r.WithContext($ctx); } bind { r as "r"; ctx as "ctx"; }
}
