arch x86_64;
platform linux, darwin, windows, freebsd;

pattern errors_Unwrap {
    name: "errors.Unwrap";
    library: "go-stdlib";
    description: "Unwraps an error";

    instr match {
        CALL errors_Unwrap
    }

    gen {
        errors.Unwrap($err);
    }

    bind {
        err as "err";
    }
}
