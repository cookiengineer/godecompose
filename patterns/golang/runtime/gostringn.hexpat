arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_gostringn {
    name: "runtime.gostringn";
    library: "go-runtime";
    description: "Creates a Go string from n bytes";

    instr match {
        MOVQ ptr, AX
        MOVQ len, BX
        CALL runtime_gostringn
    }

    gen {
        string($ptr[:$len]);
    }

    bind {
        ptr as "ptr";
        len as "len";
    }
}
