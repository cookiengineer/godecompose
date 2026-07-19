arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_makeslice {
    name: "runtime.makeslice";
    library: "go-runtime";
    description: "Creates a new slice";

    instr match {
        MOVQ typ, AX
        MOVQ len, BX
        MOVQ cap, CX
        CALL runtime_makeslice
    }

    gen {
        make([]$typ, $len, $cap);
    }

    bind {
        typ as "type";
        len as "len";
        cap as "cap";
    }
}
