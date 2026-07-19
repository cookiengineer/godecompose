arch x86_64;
platform linux, darwin, windows, freebsd;

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
