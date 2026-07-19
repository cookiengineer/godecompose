arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_makechan {
    name: "runtime.makechan";
    library: "go-runtime";
    description: "Creates a new channel";

    instr match {
        MOVQ typ, AX
        MOVQ hint, BX
        CALL runtime_makechan
    }

    gen {
        make(chan type, $hint);
    }

    bind {
        typ as "type";
        hint as "size";
    }
}
