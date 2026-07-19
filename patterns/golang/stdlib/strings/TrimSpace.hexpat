arch x86_64;
platform linux, darwin, windows, freebsd;

pattern strings_TrimSpace {
    name: "strings.TrimSpace";
    library: "go-stdlib"; instr match { CALL strings_TrimSpace }
    gen { strings.TrimSpace($s); }
    bind { s as "s"; }
}
