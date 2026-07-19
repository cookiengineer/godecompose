arch x86_64;
platform linux, darwin, windows, freebsd;

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
