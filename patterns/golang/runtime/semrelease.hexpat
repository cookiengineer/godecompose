arch x86_64;
platform linux, darwin, windows, freebsd;

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
