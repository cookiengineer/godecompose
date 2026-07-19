arch x86_64;
platform linux, darwin, windows, freebsd;

pattern errors_Is {
    name: "errors.Is";
    library: "go-stdlib";
    description: "Checks if an error wraps a target error";

    instr match {
        CALL errors_Is
    }

    gen {
        errors.Is($err, $target);
    }

    bind {
        err as "err";
        target as "target";
    }
}
