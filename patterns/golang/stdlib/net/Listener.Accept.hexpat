arch x86_64;
platform linux, darwin, windows, freebsd;

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
