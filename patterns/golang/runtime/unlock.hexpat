arch x86_64;
platform linux, darwin, windows, freebsd;

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
