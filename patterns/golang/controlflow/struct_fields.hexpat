arch x86_64;
platform linux, darwin, windows, freebsd;

pattern go_gc_write_barrier_check {
    name: "go: GC write barrier check (cmp [runtime.writeBarrier], 0)";
    library: "go-runtime";
    description: "Write barrier enabled check: cmp $0 against runtime.writeBarrier global";
    instr match {
        CMP off(base), $0
        JEQ @skip
    }
    gen { /* GC: write barrier disabled, skip */ }
    bind { off as "off"; base as "wb" }
}

pattern go_struct_field_write {
    name: "go: struct field write MOV [base+off], val";
    library: "go-controlflow";
    description: "Write to struct/stack via memory operand";
    instr match { MOV off(base), val }
    gen { /* $base+0x$off = $val */ }
    bind { val as "val"; off as "off"; base as "rcv" }
}

pattern go_struct_field_read {
    name: "go: struct field read MOV val, [base+off]";
    library: "go-controlflow";
    description: "Read from struct/stack via memory operand";
    instr match { MOV val, off(base) }
    gen { /* $val = $base+0x$off */ }
    bind { val as "val"; off as "off"; base as "rcv" }
}
