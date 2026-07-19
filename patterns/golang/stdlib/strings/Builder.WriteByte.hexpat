arch x86_64;
platform linux, darwin, windows, freebsd;

pattern strings_Builder_WriteByte {
    name: "strings.Builder.WriteByte";
    library: "go-stdlib"; instr match { CALL strings___Builder__WriteByte }
    gen { $b.WriteByte($c); }
    bind { b as "b"; c as "c"; }
}
