arch x86_64;
platform linux, darwin, windows, freebsd;

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
