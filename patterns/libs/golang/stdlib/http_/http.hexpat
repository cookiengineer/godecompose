arch x86_64;
platform linux, darwin, windows, freebsd;

pattern http_Get {
    name: "http.Get"; library: "go-stdlib";
    instr match { CALL net_http_Get }
    gen { http.Get($url); } bind { url as "url"; }
}
pattern http_Post {
    name: "http.Post"; library: "go-stdlib";
    instr match { CALL net_http_Post }
    gen { http.Post($url, $contentType, $body); } bind { url as "url"; contentType as "contentType"; body as "body"; }
}
pattern http_NewRequest {
    name: "http.NewRequest"; library: "go-stdlib";
    instr match { CALL net_http_NewRequest }
    gen { http.NewRequest($method, $url, $body); } bind { method as "method"; url as "url"; body as "body"; }
}
pattern http_ListenAndServe {
    name: "http.ListenAndServe"; library: "go-stdlib";
    instr match { CALL net_http_ListenAndServe }
    gen { http.ListenAndServe($addr, $handler); } bind { addr as "addr"; handler as "handler"; }
}
pattern http_NewServeMux {
    name: "http.NewServeMux"; library: "go-stdlib";
    instr match { CALL net_http_NewServeMux }
    gen { http.NewServeMux(); }
}
pattern http_ServeMux_HandleFunc {
    name: "http.ServeMux.HandleFunc"; library: "go-stdlib";
    instr match { CALL net_http___ServeMux__HandleFunc }
    gen { $mux.HandleFunc($pattern, $handler); } bind { mux as "mux"; pattern as "pattern"; handler as "handler"; }
}
pattern http_Client_Do {
    name: "http.Client.Do"; library: "go-stdlib";
    instr match { CALL net_http___Client__Do }
    gen { $client.Do($req); } bind { client as "client"; req as "req"; }
}
pattern http_Response_Body_Close {
    name: "http.Response.Body.Close"; library: "go-stdlib";
    instr match { CALL net_http___Response__Body_Close }
    gen { $resp.Body.Close(); } bind { resp as "resp"; }
}

pattern url_Parse {
    name: "url.Parse"; library: "go-stdlib";
    instr match { CALL net_url_Parse }
    gen { url.Parse($rawURL); } bind { rawURL as "rawURL"; }
}
pattern url_Values_Get {
    name: "url.Values.Get"; library: "go-stdlib";
    instr match { CALL net_url___Values__Get }
    gen { $v.Get($key); } bind { v as "v"; key as "key"; }
}
pattern url_Values_Set {
    name: "url.Values.Set"; library: "go-stdlib";
    instr match { CALL net_url___Values__Set }
    gen { $v.Set($key, $value); } bind { v as "v"; key as "key"; value as "value"; }
}
pattern url_Values_Add {
    name: "url.Values.Add"; library: "go-stdlib";
    instr match { CALL net_url___Values__Add }
    gen { $v.Add($key, $value); } bind { v as "v"; key as "key"; value as "value"; }
}
pattern url_Values_Encode {
    name: "url.Values.Encode"; library: "go-stdlib";
    instr match { CALL net_url___Values__Encode }
    gen { $v.Encode(); } bind { v as "v"; }
}

pattern tls_Dial {
    name: "tls.Dial"; library: "go-stdlib";
    instr match { CALL crypto_tls_Dial }
    gen { tls.Dial($network, $addr, $config); } bind { network as "network"; addr as "addr"; config as "config"; }
}
pattern tls_Listen {
    name: "tls.Listen"; library: "go-stdlib";
    instr match { CALL crypto_tls_Listen }
    gen { tls.Listen($network, $laddr, $config); } bind { network as "network"; laddr as "laddr"; config as "config"; }
}
pattern tls_LoadX509KeyPair {
    name: "tls.LoadX509KeyPair"; library: "go-stdlib";
    instr match { CALL crypto_tls_LoadX509KeyPair }
    gen { tls.LoadX509KeyPair($certFile, $keyFile); } bind { certFile as "certFile"; keyFile as "keyFile"; }
}

pattern net_Dial {
    name: "net.Dial"; library: "go-stdlib";
    instr match { CALL net_Dial }
    gen { net.Dial($network, $address); } bind { network as "network"; address as "address"; }
}
pattern net_DialTimeout {
    name: "net.DialTimeout"; library: "go-stdlib";
    instr match { CALL net_DialTimeout }
    gen { net.DialTimeout($network, $address, $timeout); } bind { network as "network"; address as "address"; timeout as "timeout"; }
}
pattern net_Listen {
    name: "net.Listen"; library: "go-stdlib";
    instr match { CALL net_Listen }
    gen { net.Listen($network, $address); } bind { network as "network"; address as "address"; }
}
pattern net_ResolveTCPAddr {
    name: "net.ResolveTCPAddr"; library: "go-stdlib";
    instr match { CALL net_ResolveTCPAddr }
    gen { net.ResolveTCPAddr($network, $address); } bind { network as "network"; address as "address"; }
}
pattern net_SplitHostPort {
    name: "net.SplitHostPort"; library: "go-stdlib";
    instr match { CALL net_SplitHostPort }
    gen { net.SplitHostPort($hostport); } bind { hostport as "hostport"; }
}
pattern net_JoinHostPort {
    name: "net.JoinHostPort"; library: "go-stdlib";
    instr match { CALL net_JoinHostPort }
    gen { net.JoinHostPort($host, $port); } bind { host as "host"; port as "port"; }
}
