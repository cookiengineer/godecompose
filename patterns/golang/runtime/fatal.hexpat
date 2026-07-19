arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_fatal {
    name: "runtime.fatal";
    library: "go-runtime";
    description: "Fatal error, exits the program";

    instr match {
        MOVQ msg, AX
        CALL runtime_fatal
    }

    gen {
        panic($msg);
    }

    bind {
        msg as "msg";
    }
}
