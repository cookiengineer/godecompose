arch x86_64;
platform linux, darwin, windows, freebsd;

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
