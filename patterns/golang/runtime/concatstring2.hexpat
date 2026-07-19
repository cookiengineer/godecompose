arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_concatstring2 {
    name: "runtime.concatstring2";
    library: "go-runtime";
    description: "Concatenates 2 strings";

    instr match {
        LEAQ buf, AX
        MOVQ str1, BX
        MOVQ str2, CX
        CALL runtime_concatstring2
    }

    gen {
        $str1 + $str2;
    }

    bind {
        str1 as "a";
        str2 as "b";
    }
}
