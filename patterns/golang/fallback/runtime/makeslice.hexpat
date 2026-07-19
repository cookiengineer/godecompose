arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_makeslice_call {
    name: "runtime.makeslice";
    library: "go-runtime";
    description: "runtime.makeslice call — slice allocation";

    instr match {
        CALL runtime_makeslice
    }

    gen {
        make([]T, $len, $cap);
    }

    bind {
        len as "len";
        cap as "cap";
    }
}
