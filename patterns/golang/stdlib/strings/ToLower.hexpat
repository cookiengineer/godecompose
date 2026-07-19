arch x86_64;
platform linux, darwin, windows, freebsd;

pattern strings_ToLower {
    name: "strings.ToLower";
    library: "go-stdlib"; instr match { CALL strings_ToLower }
    gen { strings.ToLower($s); }
    bind { s as "s"; }
}
