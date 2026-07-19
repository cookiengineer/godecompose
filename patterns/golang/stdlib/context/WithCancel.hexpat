arch x86_64;
platform linux, darwin, windows, freebsd;

pattern context_WithCancel {
    name: "context.WithCancel";
    library: "go-stdlib";
    description: "Returns a copy of parent with a new Done channel";

    instr match {
        CALL context_WithCancel
    }

    gen {
        ctx, cancel := context.WithCancel($parent);
    }

    bind {
        parent as "parent";
    }
}
