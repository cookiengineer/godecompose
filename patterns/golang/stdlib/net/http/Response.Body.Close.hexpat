arch x86_64;
platform linux, darwin, windows, freebsd;

pattern http_Response_Body_Close {
    name: "http.Response.Body.Close"; library: "go-stdlib";
    instr match { CALL net_http___Response__Body_Close }
    gen { $resp.Body.Close(); } bind { resp as "resp"; }
}
