arch x86_64;
platform linux, darwin, windows, freebsd;

pattern strings_LastIndex {
    name: "strings.LastIndex";
    library: "go-stdlib"; instr match { CALL strings_LastIndex }
    gen { strings.LastIndex($s, $substr); }
    bind { s as "s"; substr as "substr"; }
}
