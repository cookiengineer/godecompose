arch x86_64;
platform linux, darwin, windows, freebsd;

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

pattern go_defer_call {
    name: "go: defer fn()";
    library: "go-controlflow";
    description: "defer function call via runtime.deferproc";
    instr match { CALL runtime_deferproc }
    gen { defer fn() }
}

pattern go_deferstack_deferparam {
    name: "go: deferprocStack fn()";
    library: "go-controlflow";
    description: "defer function call via runtime.deferprocStack";
    instr match { CALL runtime_deferprocStack }
    gen { defer fn() }
}

pattern go_deferreturn_call {
    name: "go: deferreturn";
    library: "go-controlflow";
    description: "runtime.deferreturn — deferred return chain";
    instr match { CALL runtime_deferreturn }
    gen { /* deferred calls run */ }
}

pattern go_goroutine_go {
    name: "go: go fn() goroutine";
    library: "go-controlflow";
    description: "go fn() — starts goroutine via runtime.newproc";
    instr match { CALL runtime_newproc }
    gen { go fn() }
}

pattern go_map_lookup_ok {
    name: "go: v, ok := m[key] map access";
    library: "go-controlflow";
    description: "map access with comma-ok via runtime.mapaccess2";
    instr match { CALL runtime_mapaccess2_faststr }
    gen { v, ok := m[key] }
}

pattern go_map_lookup {
    name: "go: v := m[key] map access";
    library: "go-controlflow";
    description: "map access via runtime.mapaccess1";
    instr match { CALL runtime_mapaccess1_faststr }
    gen { v := m[key] }
}

pattern go_map_assign {
    name: "go: m[key] = val map assign";
    library: "go-controlflow";
    description: "map assignment via runtime.mapassign_faststr";
    instr match { CALL runtime_mapassign_faststr }
    gen { m[key] = val }
}

pattern go_make_chan {
    name: "go: make(chan T)";
    library: "go-controlflow";
    description: "channel creation via runtime.makechan";
    instr match { CALL runtime_makechan }
    gen { ch := make(chan T) }
}

pattern go_chan_send {
    name: "go: ch <- val channel send";
    library: "go-controlflow";
    description: "channel send via runtime.chansend1";
    instr match { CALL runtime_chansend1 }
    gen { ch <- val }
}

pattern go_chan_recv {
    name: "go: val := <-ch channel receive";
    library: "go-controlflow";
    description: "channel receive via runtime.chanrecv1";
    instr match { CALL runtime_chanrecv1 }
    gen { v := <-ch }
}

pattern go_goexit {
    name: "go: runtime.goexit";
    library: "go-controlflow";
    description: "runtime.goexit — goroutine exit";
    instr match { CALL runtime_goexit }
    gen { return }
}
