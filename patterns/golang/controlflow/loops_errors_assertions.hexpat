arch x86_64;
platform linux, darwin, windows, freebsd;

pattern go_for_loop_counted {
    name: "go: for loop counted (i=0; i<n; i++)";
    library: "go-controlflow";
    description: "Counted for loop: init zero, CMP+JGE/JLT check with INCQ increment";
    instr match {
        XOR i, i
        JMP @check
    @body:
        INC i
    @check:
        CMP i, n
        JLT @body
    }
    gen { for i := 0; i < $n; i++ { /* loop body */ } }
    bind { i as "i"; n as "n" }
}

pattern go_for_loop_range_slice {
    name: "go: for range over slice";
    library: "go-controlflow";
    description: "Range loop over slice: len check + index increment";
    instr match {
        CMP idx, len
        JGE @done
    @body:
        INC idx
        JMP @check
    @check:
    }
    gen { for idx := 0; idx < $len; idx++ { /* range body */ } }
    bind { idx as "i"; len as "n" }
}

pattern go_error_return_pattern {
    name: "go: err,val := fn(); if err != nil { return err }";
    library: "go-controlflow";
    description: "Call returning (val, err); TEST err; JEQ ok; return err";
    instr match {
        CALL fn
        TEST er, er
        JEQ @ok
        MOV er, AX
        MOV kind, BX
        RET
    @ok:
    }
    gen { val, err := $fn(); if err != nil { return err } }
    bind { fn as "fn"; er as "err"; kind as "kind" }
}

pattern go_error_return_nil_err {
    name: "go: fn returns (T, error); if err != nil { return nil, err }";
    library: "go-controlflow";
    description: "Two-return fn; TEST err; return zero,err on failure";
    instr match {
        CALL fn
        TEST er, er
        JEQ @ok
        MOV er, AX
        MOV kind, BX
        XOR val, val
        RET
    @ok:
    }
    gen { val, err := $fn(); if err != nil { return nil, err } }
    bind { fn as "fn"; er as "err"; kind as "kind"; val as "val" }
}

pattern go_type_assert_panic {
    name: "go: v := x.(T) type assertion (panics)";
    library: "go-controlflow";
    description: "Type assertion via runtime.assertE2I / assertI2I";
    instr match { CALL runtime_assertE2I }
    gen { v := x.(T) }
}

pattern go_type_assert_comma_ok {
    name: "go: v, ok := x.(T) type assertion (comma-ok)";
    library: "go-controlflow";
    description: "Type assertion with comma-ok via runtime.assertE2I2";
    instr match { CALL runtime_assertE2I2 }
    gen { v, ok := x.(T) }
}

pattern go_type_assert_iface {
    name: "go: v := x.(T) interface assertion";
    library: "go-controlflow";
    description: "Interface type assertion via runtime.assertI2I";
    instr match { CALL runtime_assertI2I }
    gen { v := x.(T) }
}

pattern go_type_switch_assert {
    name: "go: type switch assertion pattern";
    library: "go-controlflow";
    description: "Type switch via runtime.typeAssert (may panic on direct call)";
    instr match { CALL runtime_typeAssert }
    gen { v := x.(T) }
}
