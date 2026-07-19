arch x86_64;
platform linux, darwin, windows, freebsd;

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
