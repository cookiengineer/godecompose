arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_newproc {
    name: "runtime.newproc";
    library: "go-runtime";
    description: "Creates a new goroutine";

    instr match {
        MOVQ fn_ptr, AX
        CALL runtime_newproc
    }

    gen {
        go $fn_ptr;
    }

    bind {
        fn_ptr as "fn";
    }
}

pattern runtime_gopark {
    name: "runtime.gopark";
    library: "go-runtime";
    description: "Parks the current goroutine";

    instr match {
        MOVQ unlockf, AX
        MOVQ lock_ptr, BX
        MOVQ reason, CX
        CALL runtime_gopark
    }

    gen {
        gopark($unlockf, $lock_ptr, $reason);
    }

    bind {
        unlockf as "unlockf";
        lock_ptr as "lock";
        reason as "reason";
    }
}

pattern runtime_goready {
    name: "runtime.goready";
    library: "go-runtime";
    description: "Marks a goroutine as runnable";

    instr match {
        MOVQ gp, AX
        CALL runtime_goready
    }

    gen {
        goready($gp);
    }

    bind {
        gp as "gp";
    }
}

pattern runtime_goexit {
    name: "runtime.goexit";
    library: "go-runtime";
    description: "Terminates the current goroutine";

    instr match {
        MOVL zero, AX
        CALL runtime_goexit1
    }

    gen {
        runtime_goexit();
    }
}

pattern runtime_gosched {
    name: "runtime.goschedImpl";
    library: "go-runtime";
    description: "Yields the processor";

    instr match {
        CALL runtime_goschedImpl
    }

    gen {
        runtime.Gosched();
    }
}

pattern runtime_raceread {
    name: "runtime.raceread";
    library: "go-runtime";
    description: "Race detector read barrier";

    instr match {
        MOVQ addr, AX
        CALL runtime_raceread
    }

    gen {
        raceReadRange($addr, $size);
    }

    bind {
        addr as "addr";
    }
}
