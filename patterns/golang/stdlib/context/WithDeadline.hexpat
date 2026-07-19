arch x86_64;
platform linux, darwin, windows, freebsd;

pattern context_WithDeadline {
    name: "context.WithDeadline";
    library: "go-stdlib";
    description: "Returns a copy of parent with deadline";

    instr match {
        CALL context_WithDeadline
    }

    gen {
        ctx, cancel := context.WithDeadline($parent, $deadline);
    }

    bind {
        parent as "parent";
        deadline as "deadline";
    }
}
