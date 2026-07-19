arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_scanobject {
    name: "runtime.scanobject";
    library: "go-runtime";
    description: "Scans an object for GC pointers";

    instr match {
        MOVQ obj, AX
        MOVQ gcw, BX
        CALL runtime_scanobject
    }

    gen {
        scanobject($obj, $gcw);
    }

    bind {
        obj as "obj";
        gcw as "gcw";
    }
}
