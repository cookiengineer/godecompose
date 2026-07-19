arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_memhash {
    name: "runtime.memhash";
    library: "go-runtime";
    description: "Hashes n bytes of memory";

    instr match {
        MOVQ ptr, AX
        MOVQ len, BX
        MOVQ seed, CX
        CALL aeshashbody
    }

    gen {
        memhash($ptr, $seed, $len);
    }

    bind {
        ptr as "ptr";
        len as "size";
        seed as "seed";
    }
}
