arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_greyobject {
    name: "runtime.greyobject";
    library: "go-runtime";
    description: "Marks an object grey during GC";

    instr match {
        MOVQ obj, AX
        MOVQ base, BX
        MOVQ off, CX
        MOVQ gcw, DI
        CALL runtime_greyobject
    }

    gen {
        greyobject($obj, $base, $off, $gcw);
    }

    bind {
        obj as "obj";
        base as "base";
        off as "off";
        gcw as "gcw";
    }
}
