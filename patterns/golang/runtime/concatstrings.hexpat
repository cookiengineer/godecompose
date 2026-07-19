arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_concatstrings {
    name: "runtime.concatstrings";
    library: "go-runtime";
    description: "Concatenates n strings";

    instr match {
        LEAQ buf, AX
        MOVQ strs, BX
        MOVQ n, CX
        CALL runtime_concatstrings
    }

    gen {
        strings.Join($strs, "");
    }

    bind {
        strs as "strs";
        n as "n";
    }
}
