arch x86_64;
platform linux, darwin, windows, freebsd;

pattern strings_Builder_WriteRune {
    name: "strings.Builder.WriteRune";
    library: "go-stdlib"; instr match { CALL strings___Builder__WriteRune }
    gen { $b.WriteRune($r); }
    bind { b as "b"; r as "r"; }
}
