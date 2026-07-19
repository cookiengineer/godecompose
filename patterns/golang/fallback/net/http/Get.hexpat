arch x86_64;
platform linux, darwin, windows, freebsd;

pattern net_http_Get_call {
    name: "net/http.Get";
    library: "go-stdlib";
    description: "http.Get call";

    instr match {
        CALL http_Get
    }

    gen {
        resp, err := http.Get($url);
    }

    bind {
        url as "url";
    }
}
