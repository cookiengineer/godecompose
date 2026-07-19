arch x86_64;
platform linux, darwin, windows, freebsd;

pattern go_struct_field_write {
    name: "go: struct field write MOV [base+off], val";
    library: "go-controlflow";
    description: "Write to struct field via memory operand — shows offset and register";
    instr match { MOV off(base), val }
    gen { /* field write: $base+0x$off = $val */ }
    bind { val as "val"; off as "off"; base as "rcv" }
}

pattern go_struct_field_write_qw {
    name: "go: struct field write qword [base+off], val";
    library: "go-controlflow";
    description: "Write qword to struct or stack slot";
    instr match { MOVQ off(base), val }
    gen { /* mem write: $base+0x$off = $val */ }
    bind { val as "val"; off as "off"; base as "rcv" }
}

pattern go_struct_field_read {
    name: "go: struct field read MOV val, [base+off]";
    library: "go-controlflow";
    description: "Read from struct field via memory operand";
    instr match { MOV val, off(base) }
    gen { /* mem read: $val = $base+0x$off */ }
    bind { val as "val"; off as "off"; base as "rcv" }
}

pattern go_struct_field_read_qw {
    name: "go: struct field read MOVQ val, [base+off]";
    library: "go-controlflow";
    description: "Read qword from struct or stack slot";
    instr match { MOVQ val, off(base) }
    gen { /* mem read: $val = $base+0x$off */ }
    bind { val as "val"; off as "off"; base as "rcv" }
}

pattern go_struct_byte_test {
    name: "go: struct byte test CMP byte ptr [base+off], 0";
    library: "go-controlflow";
    description: "Compare byte at memory offset against zero";
    instr match {
        CMP off(base), $0
        JEQ @lab
    }
    gen { if mem[$base+0x$off] == 0 { goto $lab } }
    bind { off as "off"; base as "rcv" }
}
