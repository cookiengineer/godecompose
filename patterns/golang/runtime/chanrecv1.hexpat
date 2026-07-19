arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_chanrecv1 {
    name: "runtime.chanrecv1";
    library: "go-runtime";
    description: "Receives a value from a channel (non-blocking)";

    instr match {
        MOVQ ch, AX
        MOVQ elem, BX
        CALL runtime_chanrecv1
    }

    gen {
        // val := <-ch
    }

    bind {
        ch as "ch";
        elem as "val";
    }
}
