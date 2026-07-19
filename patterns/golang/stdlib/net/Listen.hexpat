arch x86_64;
platform linux, darwin, windows, freebsd;

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
