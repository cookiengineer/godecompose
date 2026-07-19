arch x86_64;
platform linux, darwin, windows, freebsd;

pattern context_Background {
    name: "context.Background";
    library: "go-stdlib";
    description: "Returns the background context";

    instr match {
        CALL context_Background
    }

    gen {
        ctx := context.Background();
    }
}
