arch x86_64;
platform linux, darwin, windows, freebsd;

pattern strings_Builder_WriteString {
    name: "strings.Builder.WriteString";
    library: "go-stdlib"; instr match { CALL strings___Builder__WriteString }
    gen { $b.WriteString($s); }
    bind { b as "b"; s as "s"; }
}
