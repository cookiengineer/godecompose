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
