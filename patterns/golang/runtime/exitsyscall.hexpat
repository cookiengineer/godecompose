arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_exitsyscall {
    name: "runtime.exitsyscall";
    library: "go-runtime";
    description: "Marks goroutine as exiting syscall";

    instr match {
        CALL runtime_exitsyscall
    }

    gen {
        // exitsyscall
    }
}
