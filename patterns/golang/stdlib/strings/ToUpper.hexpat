arch x86_64;
platform linux, darwin, windows, freebsd;

pattern strings_ToUpper {
    name: "strings.ToUpper";
    library: "go-stdlib"; instr match { CALL strings_ToUpper }
    gen { strings.ToUpper($s); }
    bind { s as "s"; }
}
