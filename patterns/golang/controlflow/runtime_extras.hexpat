arch x86_64;
platform linux, darwin, windows, freebsd;

pattern go_runtime_morestack_abi0 {
    name: "go: runtime.morestack.abi0";
    library: "go-runtime";
    description: "Stack growth preamble — function prologue";
    instr match { CALL runtime_morestack_abi0 }
    gen { /* stack growth check */ }
}

pattern go_runtime_morestack_noctxt_abi0 {
    name: "go: runtime.morestack_noctxt.abi0";
    library: "go-runtime";
    description: "Stack growth preamble (no closure context)";
    instr match { CALL runtime_morestack_noctxt_abi0 }
    gen { /* stack growth check (noctxt) */ }
}

pattern go_runtime_panicBounds {
    name: "go: runtime.panicBounds";
    library: "go-runtime";
    description: "Index/slice bounds check failure";
    instr match { CALL runtime_panicBounds }
    gen { panic("index out of range") }
}

pattern go_runtime_panicIndex {
    name: "go: runtime.panicIndex";
    library: "go-runtime";
    description: "Index out of range panic";
    instr match { CALL runtime_panicIndex }
    gen { panic("index out of range") }
}

pattern go_runtime_wbMove {
    name: "go: runtime.wbMove";
    library: "go-runtime";
    description: "Write barrier memory move";
    instr match { CALL runtime_wbMove }
    gen { /* wbMove GC typed memmove */ }
}

pattern go_runtime_convT {
    name: "go: runtime.convT";
    library: "go-runtime";
    description: "Convert value to interface type";
    instr match { CALL runtime_convT }
    gen { interface{}(v) }
}

pattern go_runtime_convTslice {
    name: "go: runtime.convTslice";
    library: "go-runtime";
    description: "Convert slice to interface";
    instr match { CALL runtime_convTslice }
    gen { interface{}(s) }
}

pattern go_runtime_concatstring2 {
    name: "go: runtime.concatstring2";
    library: "go-runtime";
    description: "Concatenate two strings via runtime.concatstring2";
    instr match { CALL runtime_concatstring2 }
    gen { a + b }
}

pattern go_runtime_concatstring3 {
    name: "go: runtime.concatstring3";
    library: "go-runtime";
    description: "Concatenate three strings via runtime.concatstring3";
    instr match { CALL runtime_concatstring3 }
    gen { a + b + c }
}

pattern go_runtime_concatstrings {
    name: "go: runtime.concatstrings";
    library: "go-runtime";
    description: "Concatenate strings via runtime.concatstrings";
    instr match { CALL runtime_concatstrings }
    gen { strings.Join(parts, "") }
}

pattern go_runtime_assertE2I {
    name: "go: runtime.assertE2I";
    library: "go-runtime";
    description: "Interface type assertion eface->iface";
    instr match { CALL runtime_assertE2I }
    gen { v.(T) }
}

pattern go_runtime_assertE2I2 {
    name: "go: runtime.assertE2I2";
    library: "go-runtime";
    description: "Interface type assertion with comma-ok";
    instr match { CALL runtime_assertE2I2 }
    gen { v, ok := x.(T) }
}

pattern go_runtime_assertI2I {
    name: "go: runtime.assertI2I";
    library: "go-runtime";
    description: "Interface type assertion iface->iface";
    instr match { CALL runtime_assertI2I }
    gen { v.(T) }
}

pattern go_runtime_ifaceeq {
    name: "go: runtime.ifaceeq";
    library: "go-runtime";
    description: "Interface equality comparison";
    instr match { CALL runtime_ifaceeq }
    gen { a == b }
}

pattern go_runtime_efaceeq {
    name: "go: runtime.efaceeq";
    library: "go-runtime";
    description: "Empty interface equality comparison";
    instr match { CALL runtime_efaceeq }
    gen { a == b }
}

pattern go_runtime_typeAssert {
    name: "go: runtime.typeAssert";
    library: "go-runtime";
    description: "Type assertion check (may panic)";
    instr match { CALL runtime_typeAssert }
    gen { v.(T) }
}

pattern go_strings_genSplit {
    name: "go: strings.genSplit";
    library: "go-stdlib";
    description: "Internal string split via strings.genSplit";
    instr match { CALL strings_genSplit }
    gen { strings.Split(s, sep) }
}

pattern go_strconv_ParseUint {
    name: "go: strconv.ParseUint";
    library: "go-stdlib";
    description: "Parse unsigned integer via strconv.ParseUint";
    instr match { CALL strconv_ParseUint }
    gen { strconv.ParseUint(s, base, bits) }
}

pattern go_os_OpenFile {
    name: "go: os.OpenFile";
    library: "go-stdlib";
    description: "Open file via os.OpenFile";
    instr match { CALL os_OpenFile }
    gen { os.OpenFile(name, flag, perm) }
}

pattern go_os_ReadFile {
    name: "go: os.ReadFile";
    library: "go-stdlib";
    description: "Read file via os.ReadFile";
    instr match { CALL os_ReadFile }
    gen { os.ReadFile(name) }
}

pattern go_os_WriteFile {
    name: "go: os.WriteFile";
    library: "go-stdlib";
    description: "Write file via os.WriteFile";
    instr match { CALL os_WriteFile }
    gen { os.WriteFile(name, data, perm) }
}

pattern go_net_textproto_CanonicalMIMEHeaderKey {
    name: "go: net/textproto.CanonicalMIMEHeaderKey";
    library: "go-stdlib";
    description: "Canonicalize MIME header key";
    instr match { CALL net_textproto_CanonicalMIMEHeaderKey }
    gen { textproto.CanonicalMIMEHeaderKey(k) }
}
