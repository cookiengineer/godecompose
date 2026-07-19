arch x86_64;
platform linux, darwin, windows, freebsd;

pattern net_HTTP_Get {
    name: "net/http.Get";
    library: "go-stdlib";
    description: "HTTP GET request";

    instr match {
        MOVQ url, AX
        CALL http_Get
    }

    gen {
        http.Get($url);
    }

    bind {
        url as "url";
    }
}
