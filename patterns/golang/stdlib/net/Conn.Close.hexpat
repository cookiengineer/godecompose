arch x86_64;
platform linux, darwin, windows, freebsd;

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
