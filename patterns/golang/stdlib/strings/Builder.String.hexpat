arch x86_64;
platform linux, darwin, windows, freebsd;

pattern strings_Builder_String {
    name: "strings.Builder.String";
    library: "go-stdlib"; instr match { CALL strings___Builder__String }
    gen { $b.String(); }
    bind { b as "b"; }
}
