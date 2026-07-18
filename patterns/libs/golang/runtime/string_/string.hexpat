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

pattern runtime_concatstrings {
    name: "runtime.concatstrings";
    library: "go-runtime";
    description: "Concatenates n strings";

    instr match {
        LEAQ buf, AX
        MOVQ strs, BX
        MOVQ n, CX
        CALL runtime_concatstrings
    }

    gen {
        strings.Join($strs, "");
    }

    bind {
        strs as "strs";
        n as "n";
    }
}

pattern runtime_gostring {
    name: "runtime.gostring";
    library: "go-runtime";
    description: "Creates a Go string from C string";

    instr match {
        MOVQ ptr, AX
        CALL runtime_gostring
    }

    gen {
        gostring($ptr);
    }

    bind {
        ptr as "ptr";
    }
}

pattern runtime_gostringn {
    name: "runtime.gostringn";
    library: "go-runtime";
    description: "Creates a Go string from n bytes";

    instr match {
        MOVQ ptr, AX
        MOVQ len, BX
        CALL runtime_gostringn
    }

    gen {
        string($ptr[:$len]);
    }

    bind {
        ptr as "ptr";
        len as "len";
    }
}

pattern runtime_intstring {
    name: "runtime.intstring";
    library: "go-runtime";
    description: "Converts an integer to a string";

    instr match {
        LEAQ buf, AX
        MOVQ val, BX
        CALL runtime_intstring
    }

    gen {
        string(rune($val));
    }

    bind {
        val as "r";
    }
}

pattern runtime_rawstringtmp {
    name: "runtime.rawstringtmp";
    library: "go-runtime";
    description: "Allocates a temporary string buffer";

    instr match {
        LEAQ buf, AX
        MOVQ len, BX
        CALL runtime_rawstringtmp
    }

    gen {
        // temp string buffer
    }

    bind {
        buf as "buf";
        len as "len";
    }
}
