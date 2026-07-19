arch x86_64;
platform linux, darwin, windows, freebsd;

pattern errors_New {
    name: "errors.New";
    library: "go-stdlib";
    description: "Creates a new error with the given text";

    instr match {
        CALL errors_New
    }

    gen {
        errors.New($msg);
    }

    bind {
        msg as "msg";
    }
}
