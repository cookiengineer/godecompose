arch x86_64;
platform linux, darwin, windows, freebsd;

pattern go_if_err_not_nil_log_fatal {
    name: "go: if err != nil { log.Fatal }";
    library: "go-controlflow";
    description: "if err != nil { log.Fatal(format, err) }";
    instr match {
        TEST er, er
        JEQ @after
        LEAQ fmtstr, AX
        MOVQ imm, BX
        XORL ecx, ecx
        XORL edi, edi
        MOVL esi, edi
        CALL log_Fatal
    @after:
    }
    gen { if $er != nil { log.Fatal($er) } }
    bind { er as "err"; fmtstr as "fmt"; imm as "n" }
}

pattern go_if_err_not_nil_return_err {
    name: "go: if err != nil { return err }";
    library: "go-controlflow";
    description: "if err != nil { return err }";
    instr match {
        TEST er, er
        JEQ @continue
        MOVQ er, AX
        RET
    @continue:
    }
    gen { if $er != nil { return $er } }
    bind { er as "err" }
}

pattern go_if_err_nil_return_nil_err {
    name: "go: if err != nil { return nil, err }";
    library: "go-controlflow";
    description: "if err != nil { return nil, err }";
    instr match {
        TEST er, er
        JEQ @continue
        MOVQ er, AX
        MOVQ kind, BX
        RET
    @continue:
    }
    gen { if $er != nil { return nil, $er } }
    bind { er as "err"; kind as "errKind" }
}

pattern go_if_test_jeq {
    name: "go: if v == nil (TEST+JEQ)";
    library: "go-controlflow";
    description: "nil check: TEST reg, reg; JEQ";
    instr match {
        TEST reg, reg
        JEQ @lab
    }
    gen { if $reg == nil { goto $lab } }
    bind { reg as "v" }
}

pattern go_if_test_jne {
    name: "go: if v != nil (TEST+JNE)";
    library: "go-controlflow";
    description: "nil check: TEST reg, reg; JNE";
    instr match {
        TEST reg, reg
        JNE @lab
    }
    gen { if $reg != nil { goto $lab } }
    bind { reg as "v" }
}

pattern go_cmp_imm_jne {
    name: "go: if v != const (CMP+JNE)";
    library: "go-controlflow";
    description: "Comparison: CMP reg, imm; JNE";
    instr match {
        CMP reg, imm
        JNE @lab
    }
    gen { if $reg != $imm { goto $lab } }
    bind { reg as "v"; imm as "const" }
}

pattern go_cmp_imm_jlt {
    name: "go: if v < const (CMP+JLT)";
    library: "go-controlflow";
    description: "Comparison: CMP reg, imm; JLT";
    instr match {
        CMP reg, imm
        JLT @lab
    }
    gen { if $reg < $imm { goto $lab } }
    bind { reg as "v"; imm as "const" }
}

pattern go_cmp_reg_jne {
    name: "go: if a != b (CMP reg,reg; JNE)";
    library: "go-controlflow";
    description: "Comparison: CMP reg1, reg2; JNE";
    instr match {
        CMP reg1, reg2
        JNE @lab
    }
    gen { if $reg1 != $reg2 { goto $lab } }
    bind { reg1 as "a"; reg2 as "b" }
}

pattern go_gopanic {
    name: "go: runtime.gopanic";
    library: "go-controlflow";
    description: "panic() via runtime.gopanic";
    instr match { CALL runtime_gopanic }
    gen { panic($v) }
    bind { v as "v" }
}
