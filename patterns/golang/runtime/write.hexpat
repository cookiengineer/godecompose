arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_write {
    name: "runtime.write";
    library: "go-runtime";
    description: "Writes to a file descriptor";

    instr match {
        MOVQ fd, AX
        MOVQ buf, BX
        MOVQ n, CX
        CALL runtime_write
    }

    gen {
        runtime_write($fd, $buf, $n);
    }

    bind {
        fd as "fd";
        buf as "buf";
        n as "n";
    }
}
