arch x86_64;
platform linux, darwin, windows, freebsd;

pattern http_Client_Do {
    name: "http.Client.Do"; library: "go-stdlib";
    instr match { CALL net_http___Client__Do }
    gen { $client.Do($req); } bind { client as "client"; req as "req"; }
}
