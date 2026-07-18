arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_entersyscall {
    name: "runtime.entersyscall";
    library: "go-runtime";
    description: "Marks goroutine as entering syscall";

    instr match {
        CALL runtime_entersyscall
    }

    gen {
        // entersyscall
    }
}

pattern runtime_exitsyscall {
    name: "runtime.exitsyscall";
    library: "go-runtime";
    description: "Marks goroutine as exiting syscall";

    instr match {
        CALL runtime_exitsyscall
    }

    gen {
        // exitsyscall
    }
}

pattern runtime_entersyscallblock {
    name: "runtime.entersyscallblock";
    library: "go-runtime";
    description: "Marks goroutine as entering blocking syscall";

    instr match {
        CALL runtime_entersyscallblock
    }

    gen {
        // entersyscallblock
    }
}

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

pattern runtime_open {
    name: "runtime.open";
    library: "go-runtime";
    description: "Opens a file";

    instr match {
        MOVQ name, AX
        MOVQ flags, BX
        MOVQ mode, CX
        CALL runtime_open
    }

    gen {
        runtime_open($name, $flags, $mode);
    }

    bind {
        name as "name";
        flags as "flags";
        mode as "mode";
    }
}

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

pattern runtime_futex {
    name: "runtime.futex";
    library: "go-runtime";
    description: "Linux futex syscall wrapper";

    instr match {
        MOVQ addr, AX
        MOVQ op, BX
        MOVQ val, CX
        MOVQ ts, DI
        MOVQ addr2, SI
        MOVQ val3, R8
        CALL runtime_futex
    }

    gen {
        futex($addr, $op, $val, $ts, $addr2, $val3);
    }

    bind {
        addr as "addr";
        op as "op";
        val as "val";
        ts as "ts";
        addr2 as "addr2";
        val3 as "val3";
    }
}

pattern runtime_nanotime {
    name: "runtime.nanotime1";
    library: "go-runtime";
    description: "Gets current monotonic time in nanoseconds";

    instr match {
        CALL runtime_nanotime1
    }

    gen {
        runtime_nanotime();
    }
}

pattern runtime_walltime {
    name: "runtime.walltime";
    library: "go-runtime";
    description: "Gets current wall clock time";

    instr match {
        LEAQ sec, AX
        LEAQ nsec, BX
        CALL runtime_walltime
    }

    gen {
        runtime_walltime(&$sec, &$nsec);
    }

    bind {
        sec as "sec";
        nsec as "nsec";
    }
}
