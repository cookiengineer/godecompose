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

pattern runtime_closechan {
    name: "runtime.closechan";
    library: "go-runtime";
    description: "Closes a channel";

    instr match {
        MOVQ ch, AX
        CALL runtime_closechan
    }

    gen {
        close($ch);
    }

    bind {
        ch as "ch";
    }
}

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

pattern runtime_makechan {
    name: "runtime.makechan";
    library: "go-runtime";
    description: "Creates a new channel";

    instr match {
        MOVQ typ, AX
        MOVQ hint, BX
        CALL runtime_makechan
    }

    gen {
        make(chan type, $hint);
    }

    bind {
        typ as "type";
        hint as "size";
    }
}
