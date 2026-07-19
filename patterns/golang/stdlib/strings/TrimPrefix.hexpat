arch x86_64;
platform linux, darwin, windows, freebsd;

pattern strings_TrimPrefix {
    name: "strings.TrimPrefix";
    library: "go-stdlib"; instr match { CALL strings_TrimPrefix }
    gen { strings.TrimPrefix($s, $prefix); }
    bind { s as "s"; prefix as "prefix"; }
}
