arch x86_64;
platform linux, darwin, windows, freebsd;

pattern strings_Count {
    name: "strings.Count";
    library: "go-stdlib"; instr match { CALL strings_Count }
    gen { strings.Count($s, $substr); }
    bind { s as "s"; substr as "substr"; }
}
