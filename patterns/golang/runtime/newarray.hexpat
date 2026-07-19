arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_newarray {
    name: "runtime.newarray";
    library: "go-runtime";
    description: "Allocates an array of n elements of the given type";

    instr match {
        MOVQ typ, AX
        MOVQ count, BX
        CALL runtime_mallocgc
    }

    gen {
        make([]$typ, $count);
    }

    bind {
        typ as "type";
        count as "n";
    }
}
