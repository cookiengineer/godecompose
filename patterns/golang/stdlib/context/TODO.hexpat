arch x86_64;
platform linux, darwin, windows, freebsd;

pattern context_TODO {
    name: "context.TODO";
    library: "go-stdlib";
    description: "Returns the TODO context";

    instr match {
        CALL context_TODO
    }

    gen {
        ctx := context.TODO();
    }
}
