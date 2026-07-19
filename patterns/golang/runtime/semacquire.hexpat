arch x86_64;
platform linux, darwin, windows, freebsd;

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
