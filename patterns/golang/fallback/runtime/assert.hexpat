arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_type_assert {
    name: "runtime.assert";
    library: "go-runtime";
    description: "Type assertion operation";

    instr match {
        CALL runtime_assertE2I
        | CALL runtime_assertI2I
        | CALL runtime_assertE2I2
    }

    gen {
        // type assertion: val.(T)
    }

    bind {
        val as "val";
    }
}
