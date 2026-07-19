arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_throw {
    name: "runtime.throw";
    library: "go-runtime";
    description: "Throws a fatal runtime error";

    instr match {
        MOVQ msg, AX
        CALL runtime_throw
    }

    gen {
        panic($msg);
    }

    bind {
        msg as "msg";
    }
}
