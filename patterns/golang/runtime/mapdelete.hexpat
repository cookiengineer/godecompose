arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_mapdelete {
    name: "runtime.mapdelete";
    library: "go-runtime";
    description: "Deletes a key from a map";

    instr match {
        MOVQ typ, AX
        MOVQ hmap, BX
        MOVQ key, CX
        CALL runtime_mapdelete
    }

    gen {
        delete($hmap, $key);
    }

    bind {
        typ as "type";
        hmap as "m";
        key as "key";
    }
}
