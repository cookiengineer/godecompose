arch x86_64;
platform linux, darwin, windows, freebsd;

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
