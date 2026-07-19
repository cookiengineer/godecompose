arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_goready {
    name: "runtime.goready";
    library: "go-runtime";
    description: "Marks a goroutine as runnable";

    instr match {
        MOVQ gp, AX
        CALL runtime_goready
    }

    gen {
        goready($gp);
    }

    bind {
        gp as "gp";
    }
}
