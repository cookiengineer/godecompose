arch x86_64;
platform linux, darwin, windows, freebsd;

pattern strings_HasPrefix {
    name: "strings.HasPrefix";
    library: "go-stdlib"; instr match { CALL strings_HasPrefix }
    gen { strings.HasPrefix($s, $prefix); }
    bind { s as "s"; prefix as "prefix"; }
}
