arch x86_64;
platform linux, darwin, windows, freebsd;

pattern strings_Index {
    name: "strings.Index";
    library: "go-stdlib"; instr match { CALL strings_Index }
    gen { strings.Index($s, $substr); }
    bind { s as "s"; substr as "substr"; }
}
