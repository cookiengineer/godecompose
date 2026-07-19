arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_panic {
    name: "runtime.gopanic";
    library: "go-runtime";
    description: "Panic call";

    instr match {
        CALL runtime_gopanic
    }

    gen {
        panic($val);
    }

    bind {
        val as "val";
    }
}
