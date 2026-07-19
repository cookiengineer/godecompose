arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_chan_send {
    name: "runtime.chansend";
    library: "go-runtime";
    description: "Channel send operation";

    instr match {
        CALL runtime_chansend1
    }

    gen {
        // ch <- val
    }

    bind {
        ch as "ch";
        val as "val";
    }
}
