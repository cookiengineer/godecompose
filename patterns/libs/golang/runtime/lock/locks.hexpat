arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_lock {
    name: "runtime.lock";
    library: "go-runtime";
    description: "Acquires a runtime mutex";

    instr match {
        MOVQ lk, AX
        CALL runtime_lock
    }

    gen {
        runtime_lock($lk);
    }

    bind {
        lk as "lk";
    }
}

pattern runtime_unlock {
    name: "runtime.unlock";
    library: "go-runtime";
    description: "Releases a runtime mutex";

    instr match {
        MOVQ lk, AX
        CALL runtime_unlock
    }

    gen {
        runtime_unlock($lk);
    }

    bind {
        lk as "lk";
    }
}

pattern runtime_notesleep {
    name: "runtime.notesleep";
    library: "go-runtime";
    description: "Sleeps on a runtime note";

    instr match {
        MOVQ n, AX
        CALL runtime_notesleep
    }

    gen {
        notesleep($n);
    }

    bind {
        n as "n";
    }
}

pattern runtime_notewakeup {
    name: "runtime.notewakeup";
    library: "go-runtime";
    description: "Wakes up a runtime note";

    instr match {
        MOVQ n, AX
        CALL runtime_notewakeup
    }

    gen {
        notewakeup($n);
    }

    bind {
        n as "n";
    }
}

pattern runtime_notetsleep {
    name: "runtime.notetsleep";
    library: "go-runtime";
    description: "Sleeps on a runtime note with timeout";

    instr match {
        MOVQ n, AX
        MOVQ ns, BX
        CALL runtime_notetsleep
    }

    gen {
        notetsleep($n, $ns);
    }

    bind {
        n as "n";
        ns as "ns";
    }
}

pattern runtime_semacquire {
    name: "runtime.semacquire";
    library: "go-runtime";
    description: "Acquires a semaphore";

    instr match {
        MOVQ addr, AX
        CALL runtime_semacquire
    }

    gen {
        sync_runtime_Semacquire($addr);
    }

    bind {
        addr as "addr";
    }
}

pattern runtime_semrelease {
    name: "runtime.semrelease";
    library: "go-runtime";
    description: "Releases a semaphore";

    instr match {
        MOVQ addr, AX
        CALL runtime_semrelease
    }

    gen {
        sync_runtime_Semrelease($addr);
    }

    bind {
        addr as "addr";
    }
}
