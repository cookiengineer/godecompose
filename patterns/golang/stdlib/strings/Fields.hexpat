arch x86_64;
platform linux, darwin, windows, freebsd;

pattern strings_Fields {
    name: "strings.Fields";
    library: "go-stdlib"; instr match { CALL strings_Fields }
    gen { strings.Fields($s); }
    bind { s as "s"; }
}
