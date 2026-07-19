arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_chansend1 {
    name: "runtime.chansend1";
    library: "go-runtime";
    description: "Sends a value on a channel (non-blocking)";

    instr match {
        MOVQ ch, AX
        MOVQ elem, BX
        CALL runtime_chansend1
    }

    gen {
        // ch <- val
    }

    bind {
        ch as "ch";
        elem as "val";
    }
}
