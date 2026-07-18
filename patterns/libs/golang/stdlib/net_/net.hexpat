arch x86_64;
platform linux, darwin, windows, freebsd;

pattern net_Dial {
    name: "net.Dial";
    library: "go-stdlib";
    description: "Connects to a network address";

    instr match {
        MOVQ network, AX
        MOVQ address, BX
        CALL net_Dial
    }

    gen {
        net.Dial($network, $address);
    }

    bind {
        network as "network";
        address as "address";
    }
}

pattern net_Listen {
    name: "net.Listen";
    library: "go-stdlib";
    description: "Creates a network listener";

    instr match {
        MOVQ network, AX
        MOVQ address, BX
        CALL net_Listen
    }

    gen {
        net.Listen($network, $address);
    }

    bind {
        network as "network";
        address as "address";
    }
}

pattern net_Accept {
    name: "net.Listener.Accept";
    library: "go-stdlib";
    description: "Accepts a network connection";

    instr match {
        MOVQ ln, AX
        CALL net_Listener_Accept
    }

    gen {
        $ln.Accept();
    }

    bind {
        ln as "ln";
    }
}

pattern net_Conn_Read {
    name: "net.Conn.Read";
    library: "go-stdlib";
    description: "Reads from a network connection";

    instr match {
        MOVQ conn, AX
        MOVQ buf, BX
        CALL net_Conn_Read
    }

    gen {
        $conn.Read($buf);
    }

    bind {
        conn as "conn";
        buf as "buf";
    }
}

pattern net_Conn_Write {
    name: "net.Conn.Write";
    library: "go-stdlib";
    description: "Writes to a network connection";

    instr match {
        MOVQ conn, AX
        MOVQ buf, BX
        CALL net_Conn_Write
    }

    gen {
        $conn.Write($buf);
    }

    bind {
        conn as "conn";
        buf as "buf";
    }
}

pattern net_Conn_Close {
    name: "net.Conn.Close";
    library: "go-stdlib";
    description: "Closes a network connection";

    instr match {
        MOVQ conn, AX
        CALL net_Conn_Close
    }

    gen {
        $conn.Close();
    }

    bind {
        conn as "conn";
    }
}

pattern net_ResolveTCPAddr {
    name: "net.ResolveTCPAddr";
    library: "go-stdlib";
    description: "Resolves a TCP address";

    instr match {
        MOVQ network, AX
        MOVQ address, BX
        CALL net_ResolveTCPAddr
    }

    gen {
        net.ResolveTCPAddr($network, $address);
    }

    bind {
        network as "network";
        address as "address";
    }
}

pattern net_DialTCP {
    name: "net.DialTCP";
    library: "go-stdlib";
    description: "Dials a TCP connection";

    instr match {
        MOVQ network, AX
        MOVQ laddr, BX
        MOVQ raddr, CX
        CALL net_DialTCP
    }

    gen {
        net.DialTCP($network, $laddr, $raddr);
    }

    bind {
        network as "network";
        laddr as "laddr";
        raddr as "raddr";
    }
}

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
