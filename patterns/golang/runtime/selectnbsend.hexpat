arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_selectnbsend {
    name: "runtime.selectnbsend";
    library: "go-runtime";
    description: "Non-blocking send for select";

    instr match {
        MOVQ ch, AX
        MOVQ elem, BX
        CALL runtime_selectnbsend
    }

    gen {
        // select send: ch <- val
    }

    bind {
        ch as "ch";
        elem as "val";
    }
}
