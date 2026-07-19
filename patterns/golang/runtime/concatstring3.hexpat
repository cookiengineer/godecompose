arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_concatstring3 {
    name: "runtime.concatstring3";
    library: "go-runtime";
    description: "Concatenates 3 strings";

    instr match {
        LEAQ buf, AX
        MOVQ str1, BX
        MOVQ str2, CX
        MOVQ str3, DI
        CALL runtime_concatstring3
    }

    gen {
        $str1 + $str2 + $str3;
    }

    bind {
        str1 as "a";
        str2 as "b";
        str3 as "c";
    }
}
