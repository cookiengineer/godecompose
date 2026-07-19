arch x86_64;
platform linux, darwin, windows, freebsd;

pattern go_runtime_memclr {
    name: "go: runtime.memclrNoHeapPointers";
    library: "go-runtime";
    description: "Clear memory without heap pointers";
    instr match { CALL runtime_memclrNoHeapPointers }
    gen { /* memclr */ }
}

pattern go_runtime_typedmemmove {
    name: "go: runtime.typedmemmove";
    library: "go-runtime";
    description: "Typed memory move (GC-aware copy)";
    instr match { CALL runtime_typedmemmove }
    gen { /* typed memmove */ }
}

pattern go_runtime_typedmemclr {
    name: "go: runtime.typedmemclr";
    library: "go-runtime";
    description: "Typed memory clear (GC-aware zeroing)";
    instr match { CALL runtime_typedmemclr }
    gen { /* typed memclr */ }
}

pattern go_runtime_slicecopy {
    name: "go: runtime.slicecopy";
    library: "go-runtime";
    description: "Copy slice data via runtime.slicecopy";
    instr match { CALL runtime_slicecopy }
    gen { copy(dst, src) }
}

pattern go_runtime_newarray {
    name: "go: runtime.newarray";
    library: "go-runtime";
    description: "Allocate array via runtime.newarray";
    instr match { CALL runtime_newarray }
    gen { new([n]T) }
}

pattern go_runtime_rawstringtmp {
    name: "go: runtime.rawstringtmp";
    library: "go-runtime";
    description: "Allocate uninitialized string via runtime.rawstringtmp";
    instr match { CALL runtime_rawstringtmp }
    gen { /* string allocation */ }
}

pattern go_runtime_gostring {
    name: "go: runtime.gostring";
    library: "go-runtime";
    description: "Create string from bytes via runtime.gostring";
    instr match { CALL runtime_gostring }
    gen { string(b) }
}

pattern go_runtime_gostringn {
    name: "go: runtime.gostringn";
    library: "go-runtime";
    description: "Create string from n bytes via runtime.gostringn";
    instr match { CALL runtime_gostringn }
    gen { string(b[:n]) }
}

pattern go_runtime_memhash {
    name: "go: runtime.memhash";
    library: "go-runtime";
    description: "Hash memory block via runtime.memhash";
    instr match { CALL runtime_memhash }
    gen { /* hash */ }
}

pattern go_runtime_nilinterhash {
    name: "go: runtime.nilinterhash";
    library: "go-runtime";
    description: "Hash nil interface via runtime.nilinterhash";
    instr match { CALL runtime_nilinterhash }
    gen { /* nilinter hash */ }
}

pattern go_runtime_typehash {
    name: "go: runtime.typehash";
    library: "go-runtime";
    description: "Hash typed value via runtime.typehash";
    instr match { CALL runtime_typehash }
    gen { /* type hash */ }
}

pattern go_runtime_convT2E {
    name: "go: runtime.convT2E";
    library: "go-runtime";
    description: "Convert to empty interface (eface) via runtime.convT2E";
    instr match { CALL runtime_convT2E }
    gen { interface{}(v) }
}

pattern go_runtime_convT2I {
    name: "go: runtime.convT2I";
    library: "go-runtime";
    description: "Convert to interface (iface) via runtime.convT2I";
    instr match { CALL runtime_convT2I }
    gen { T(v) }
}

pattern go_runtime_nanotime {
    name: "go: runtime.nanotime";
    library: "go-runtime";
    description: "Get current nanosecond time via runtime.nanotime";
    instr match { CALL runtime_nanotime1 }
    gen { time.Now() }
}

pattern go_runtime_walltime {
    name: "go: runtime.walltime";
    library: "go-runtime";
    description: "Get wall clock time via runtime.walltime";
    instr match { CALL runtime_walltime }
    gen { time.Now() }
}

pattern go_runtime_makemap_small {
    name: "go: runtime.makemap_small";
    library: "go-runtime";
    description: "Allocate small map via runtime.makemap_small";
    instr match { CALL runtime_makemap_small }
    gen { make(map[K]V) }
}

pattern go_uintptr_from_int {
    name: "go: uintptr(int) data conversion";
    library: "go-controlflow";
    description: "Convert int to uintptr for syscall/capabilities";
    instr match { CALL runtime_convT64 }
    gen { uintptr(v) }
}

pattern go_value_to_interface {
    name: "go: value to interface{} box";
    library: "go-controlflow";
    description: "Box a value into an empty interface";
    instr match { CALL runtime_convT2E }
    gen { interface{}(v) }
}
