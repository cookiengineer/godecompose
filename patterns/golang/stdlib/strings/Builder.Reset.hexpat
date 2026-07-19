arch x86_64;
platform linux, darwin, windows, freebsd;

pattern strings_Builder_Reset {
    name: "strings.Builder.Reset";
    library: "go-stdlib"; instr match { CALL strings___Builder__Reset }
    gen { $b.Reset(); }
    bind { b as "b"; }
}
