arch x86_64;
platform linux, darwin, windows, freebsd;

pattern net_HTTP_Post {
    name: "net/http.Post";
    library: "go-stdlib";
    description: "HTTP POST request";

    instr match {
        MOVQ url, AX
        MOVQ contentType, BX
        MOVQ body, CX
        CALL http_Post
    }

    gen {
        http.Post($url, $contentType, $body);
    }

    bind {
        url as "url";
        contentType as "contentType";
        body as "body";
    }
}
