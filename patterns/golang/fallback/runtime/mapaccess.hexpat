arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_map_access {
    name: "runtime.mapaccess";
    library: "go-runtime";
    description: "Map access operation";

    instr match {
        CALL runtime_mapaccess1
        | CALL runtime_mapaccess2
    }

    gen {
        // map access: $m[$key]
    }

    bind {
        m as "m";
        key as "key";
    }
}
