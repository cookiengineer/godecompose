arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_selectnbrecv {
    name: "runtime.selectnbrecv";
    library: "go-runtime";
    description: "Non-blocking receive for select";

    instr match {
        MOVQ ch, AX
        MOVQ elem, BX
        CALL runtime_selectnbrecv
    }

    gen {
        // select recv: <-ch
    }

    bind {
        ch as "ch";
        elem as "val";
    }
}
