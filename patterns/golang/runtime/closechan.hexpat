arch x86_64;
platform linux, darwin, windows, freebsd;

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
