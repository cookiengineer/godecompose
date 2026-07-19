arch x86_64;
platform linux, darwin, windows, freebsd;

pattern context_WithValue {
    name: "context.WithValue";
    library: "go-stdlib";
    description: "Returns a copy of parent with associated key-value pair";

    instr match {
        CALL context_WithValue
    }

    gen {
        ctx := context.WithValue($parent, $key, $val);
    }

    bind {
        parent as "parent";
        key as "key";
        val as "val";
    }
}
