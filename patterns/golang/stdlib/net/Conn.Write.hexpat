arch x86_64;
platform linux, darwin, windows, freebsd;

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
