arch x86_64;
platform linux, darwin, windows, freebsd;

pattern runtime_stringtoslicebyte {
    name: "runtime.stringtoslicebyte";
    library: "go-runtime";
    description: "Converts a string to a byte slice";

    instr match {
        MOVQ str_ptr, AX
        MOVQ str_len, BX
        CALL runtime_stringtoslicebyte
    }

    gen {
        []byte($str);
    }

    bind {
        str_ptr as "str";
        str_len as "len";
    }
}
