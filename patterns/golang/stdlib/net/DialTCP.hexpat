arch x86_64;
platform linux, darwin, windows, freebsd;

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
