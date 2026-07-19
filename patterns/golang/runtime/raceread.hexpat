arch x86_64;
platform linux, darwin, windows, freebsd;

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
