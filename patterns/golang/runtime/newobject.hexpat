arch x86_64;
platform linux, darwin, windows, freebsd;

pattern go_runtime_newobject {
    name: "runtime.newobject";
    library: "go-runtime";
    description: "Go runtime newobject — allocates a new zeroed object";

    instr match {
        MOVQ type_ptr, AX
        CALL runtime_newobject
    }

    gen {
        *$type_ptr = new(Type);
    }

    bind {
        type_ptr as "type";
    }
}
