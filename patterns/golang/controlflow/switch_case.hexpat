arch x86_64;
platform linux, darwin, windows, freebsd;

pattern go_switch_cmp_chain {
    name: "go: switch/case CMP+JEQ chain on same reg";
    library: "go-controlflow";
    description: "Switch/case detected by two sequential CMP+JEQ on same register";
    instr match {
        CMP reg, c1
        JEQ @case1
        CMP reg, c2
        JEQ @case2
    }
    gen { switch $reg { case $c1: goto $case1; case $c2: goto $case2 } }
    bind { reg as "v"; c1 as "a"; c2 as "b" }
}

pattern go_switch_threeway {
    name: "go: switch/case three-way CMP+JEQ+JMP";
    library: "go-controlflow";
    description: "Three-way switch with default: CMP+JEQ, CMP+JEQ, JMP";
    instr match {
        CMP reg, c1
        JEQ @case1
        CMP reg, c2
        JEQ @case2
        JMP @default
    }
    gen { switch $reg { case $c1: goto $case1; case $c2: goto $case2; default: goto $default } }
    bind { reg as "v"; c1 as "a"; c2 as "b" }
}
