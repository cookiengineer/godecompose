arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_closefd {
    name: "runtime.closefd";
    library: "go-runtime";
    description: "Closes a file descriptor";

    instr match {
        MOVQ fd, AX
        CALL runtime_closefd
    }

    gen {
        runtime_closefd($fd);
    }

    bind {
        fd as "fd";
    }
}
