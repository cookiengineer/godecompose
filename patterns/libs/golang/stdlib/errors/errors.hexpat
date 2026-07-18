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

pattern errors_As {
    name: "errors.As";
    library: "go-stdlib";
    description: "Finds the first error in chain matching target type";

    instr match {
        CALL errors_As
    }

    gen {
        errors.As($err, &$target);
    }

    bind {
        err as "err";
        target as "target";
    }
}

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
