arch x86_64;
platform linux, darwin, windows, freebsd;

pattern strings_Repeat {
    name: "strings.Repeat";
    library: "go-stdlib"; instr match { CALL strings_Repeat }
    gen { strings.Repeat($s, $count); }
    bind { s as "s"; count as "count"; }
}
