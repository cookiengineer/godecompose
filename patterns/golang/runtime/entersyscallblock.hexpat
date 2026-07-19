arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_entersyscallblock {
    name: "runtime.entersyscallblock";
    library: "go-runtime";
    description: "Marks goroutine as entering blocking syscall";

    instr match {
        CALL runtime_entersyscallblock
    }

    gen {
        // entersyscallblock
    }
}
