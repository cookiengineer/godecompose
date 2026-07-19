arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_reflectcall {
    name: "runtime.reflectcall";
    library: "go-runtime";
    description: "Calls a function through reflection";

    instr match {
        MOVQ fn, AX
        MOVQ arg, BX
        MOVQ argsize, CX
        MOVQ ret, DI
        MOVQ retsize, SI
        CALL runtime_reflectcall
    }

    gen {
        // reflectcall($fn, $arg, $argsize, $ret, $retsize)
    }

    bind {
        fn as "fn";
        arg as "args";
        argsize as "argsSize";
        ret as "ret";
        retsize as "retSize";
    }
}
