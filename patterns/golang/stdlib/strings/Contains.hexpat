arch x86_64;
platform linux, darwin, windows, freebsd;

pattern strings_Contains {
    name: "strings.Contains";
    library: "go-stdlib"; instr match { CALL strings_Contains }
    gen { strings.Contains($s, $substr); }
    bind { s as "s"; substr as "substr"; }
}
