arch x86_64;
platform linux, darwin, windows, freebsd;

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
