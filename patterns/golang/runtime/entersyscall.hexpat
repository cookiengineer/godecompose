arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_entersyscall {
    name: "runtime.entersyscall";
    library: "go-runtime";
    description: "Marks goroutine as entering syscall";

    instr match {
        CALL runtime_entersyscall
    }

    gen {
        // entersyscall
    }
}
