arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_mapiternext {
    name: "runtime.mapiternext";
    library: "go-runtime";
    description: "Advances a map iterator";

    instr match {
        LEAQ it_ptr, AX
        CALL runtime_mapiternext
    }

    gen {
        // next iteration
    }

    bind {
        it_ptr as "it";
    }
}
