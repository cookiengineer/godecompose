arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_futex {
    name: "runtime.futex";
    library: "go-runtime";
    description: "Linux futex syscall wrapper";

    instr match {
        MOVQ addr, AX
        MOVQ op, BX
        MOVQ val, CX
        MOVQ ts, DI
        MOVQ addr2, SI
        MOVQ val3, R8
        CALL runtime_futex
    }

    gen {
        futex($addr, $op, $val, $ts, $addr2, $val3);
    }

    bind {
        addr as "addr";
        op as "op";
        val as "val";
        ts as "ts";
        addr2 as "addr2";
        val3 as "val3";
    }
}
