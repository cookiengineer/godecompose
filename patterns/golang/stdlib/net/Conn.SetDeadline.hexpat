arch x86_64;
platform linux, darwin, windows, freebsd;

pattern net_Conn_SetDeadline {
    name: "net.Conn.SetDeadline"; library: "go-stdlib";
    instr match { CALL net_Conn_SetDeadline }
    gen { $conn.SetDeadline($t); } bind { conn as "conn"; t as "t"; }
}
