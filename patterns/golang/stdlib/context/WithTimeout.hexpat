arch x86_64;
platform linux, darwin, windows, freebsd;

pattern context_WithTimeout {
    name: "context.WithTimeout";
    library: "go-stdlib";
    description: "Returns a copy of parent with timeout";

    instr match {
        CALL context_WithTimeout
    }

    gen {
        ctx, cancel := context.WithTimeout($parent, $timeout);
    }

    bind {
        parent as "parent";
        timeout as "timeout";
    }
}
