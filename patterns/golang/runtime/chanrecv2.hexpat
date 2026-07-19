arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_chanrecv2 {
    name: "runtime.chanrecv2";
    library: "go-runtime";
    description: "Receives a value with ok flag from a channel";

    instr match {
        MOVQ ch, AX
        MOVQ elem, BX
        CALL runtime_chanrecv2
    }

    gen {
        // val, ok := <-ch
    }

    bind {
        ch as "ch";
        elem as "val";
    }
}
