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
pattern http_ListenAndServeTLS {
    name: "http.ListenAndServeTLS"; library: "go-stdlib";
    instr match { CALL net_http_ListenAndServeTLS }
    gen { http.ListenAndServeTLS($addr, $certFile, $keyFile, $handler); } bind { addr as "addr"; certFile as "certFile"; keyFile as "keyFile"; handler as "handler"; }
}
pattern http_PostForm {
    name: "http.PostForm"; library: "go-stdlib";
    instr match { CALL net_http_PostForm }
    gen { http.PostForm($url, $data); } bind { url as "url"; data as "data"; }
}
pattern http_Head {
    name: "http.Head"; library: "go-stdlib";
    instr match { CALL net_http_Head }
    gen { http.Head($url); } bind { url as "url"; }
}
pattern http_Request_WithContext {
    name: "http.Request.WithContext"; library: "go-stdlib";
    instr match { CALL net_http___Request__WithContext }
    gen { $r.WithContext($ctx); } bind { r as "r"; ctx as "ctx"; }
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
pattern url_ParseRequestURI {
    name: "url.ParseRequestURI"; library: "go-stdlib";
    instr match { CALL net_url_ParseRequestURI }
    gen { url.ParseRequestURI($rawURL); } bind { rawURL as "rawURL"; }
}
pattern url_Values_Del {
    name: "url.Values.Del"; library: "go-stdlib";
    instr match { CALL net_url___Values__Del }
    gen { $v.Del($key); } bind { v as "v"; key as "key"; }
}
pattern url_QueryEscape {
    name: "url.QueryEscape"; library: "go-stdlib";
    instr match { CALL net_url_QueryEscape }
    gen { url.QueryEscape($s); } bind { s as "s"; }
}
pattern url_QueryUnescape {
    name: "url.QueryUnescape"; library: "go-stdlib";
    instr match { CALL net_url_QueryUnescape }
    gen { url.QueryUnescape($s); } bind { s as "s"; }
}
pattern url_PathEscape {
    name: "url.PathEscape"; library: "go-stdlib";
    instr match { CALL net_url_PathEscape }
    gen { url.PathEscape($s); } bind { s as "s"; }
}
pattern url_PathUnescape {
    name: "url.PathUnescape"; library: "go-stdlib";
    instr match { CALL net_url_PathUnescape }
    gen { url.PathUnescape($s); } bind { s as "s"; }
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
pattern tls_DialWithDialer {
    name: "tls.DialWithDialer"; library: "go-stdlib";
    instr match { CALL crypto_tls_DialWithDialer }
    gen { tls.DialWithDialer($dialer, $network, $addr, $config); } bind { dialer as "dialer"; network as "network"; addr as "addr"; config as "config"; }
}
pattern tls_NewListener {
    name: "tls.NewListener"; library: "go-stdlib";
    instr match { CALL crypto_tls_NewListener }
    gen { tls.NewListener($inner, $config); } bind { inner as "inner"; config as "config"; }
}
pattern tls_X509KeyPair {
    name: "tls.X509KeyPair"; library: "go-stdlib";
    instr match { CALL crypto_tls_X509KeyPair }
    gen { tls.X509KeyPair($certPEMBlock, $keyPEMBlock); } bind { certPEMBlock as "certPEMBlock"; keyPEMBlock as "keyPEMBlock"; }
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
pattern net_ListenPacket {
    name: "net.ListenPacket"; library: "go-stdlib";
    instr match { CALL net_ListenPacket }
    gen { net.ListenPacket($network, $address); } bind { network as "network"; address as "address"; }
}
pattern net_ResolveIPAddr {
    name: "net.ResolveIPAddr"; library: "go-stdlib";
    instr match { CALL net_ResolveIPAddr }
    gen { net.ResolveIPAddr($network, $address); } bind { network as "network"; address as "address"; }
}
pattern net_ResolveUDPAddr {
    name: "net.ResolveUDPAddr"; library: "go-stdlib";
    instr match { CALL net_ResolveUDPAddr }
    gen { net.ResolveUDPAddr($network, $address); } bind { network as "network"; address as "address"; }
}
pattern net_LookupHost {
    name: "net.LookupHost"; library: "go-stdlib";
    instr match { CALL net_LookupHost }
    gen { net.LookupHost($host); } bind { host as "host"; }
}
pattern net_LookupAddr {
    name: "net.LookupAddr"; library: "go-stdlib";
    instr match { CALL net_LookupAddr }
    gen { net.LookupAddr($addr); } bind { addr as "addr"; }
}
pattern net_InterfaceAddrs {
    name: "net.InterfaceAddrs"; library: "go-stdlib";
    instr match { CALL net_InterfaceAddrs }
    gen { net.InterfaceAddrs(); }
}
pattern net_InterfaceByName {
    name: "net.InterfaceByName"; library: "go-stdlib";
    instr match { CALL net_InterfaceByName }
    gen { net.InterfaceByName($name); } bind { name as "name"; }
}
pattern net_Conn_SetDeadline {
    name: "net.Conn.SetDeadline"; library: "go-stdlib";
    instr match { CALL net_Conn_SetDeadline }
    gen { $conn.SetDeadline($t); } bind { conn as "conn"; t as "t"; }
}
