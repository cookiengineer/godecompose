arch x86_64;
platform linux, darwin, windows, freebsd;

pattern errors_Join {
    name: "errors.Join";
    library: "go-stdlib";
    description: "Joins multiple errors into a single error";

    instr match {
        CALL errors_Join
    }

    gen {
        errors.Join($errs...);
    }

    bind {
        errs as "errs";
    }
}
