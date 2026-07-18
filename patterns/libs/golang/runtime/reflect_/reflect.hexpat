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

pattern runtime_reflectcallmove {
    name: "runtime.reflectcallmove";
    library: "go-runtime";
    description: "Moves call arguments for reflection";

    instr match {
        MOVQ typ, AX
        MOVQ src, BX
        MOVQ dst, CX
        CALL runtime_reflectcallmove
    }

    gen {
        // reflectcallmove($typ, $src, $dst)
    }

    bind {
        typ as "typ";
        src as "src";
        dst as "dst";
    }
}

pattern runtime_callMethod {
    name: "reflect.callMethod";
    library: "go-runtime";
    description: "Calls a method via reflection";

    instr match {
        MOVQ method, AX
        MOVQ frame, BX
        MOVQ frametype, CX
        CALL reflect_callMethod
    }

    gen {
        // reflect.callMethod($method, $frame);
    }

    bind {
        method as "method";
        frame as "frame";
    }
}
