arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_rawstringtmp {
    name: "runtime.rawstringtmp";
    library: "go-runtime";
    description: "Allocates a temporary string buffer";

    instr match {
        LEAQ buf, AX
        MOVQ len, BX
        CALL runtime_rawstringtmp
    }

    gen {
        // temp string buffer
    }

    bind {
        buf as "buf";
        len as "len";
    }
}
