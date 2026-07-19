arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_read {
    name: "runtime.read";
    library: "go-runtime";
    description: "Reads from a file descriptor";

    instr match {
        MOVQ fd, AX
        MOVQ buf, BX
        MOVQ n, CX
        CALL runtime_read
    }

    gen {
        runtime_read($fd, $buf, $n);
    }

    bind {
        fd as "fd";
        buf as "buf";
        n as "n";
    }
}
