arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_chan_recv {
    name: "runtime.chanrecv";
    library: "go-runtime";
    description: "Channel receive operation";

    instr match {
        CALL runtime_chanrecv1
        | CALL runtime_chanrecv2
    }

    gen {
        <-$ch;
    }

    bind {
        ch as "ch";
    }
}
