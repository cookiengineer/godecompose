arch x86_64;
platform linux, darwin, windows, freebsd;

pattern strings_NewReader {
    name: "strings.NewReader";
    library: "go-stdlib"; instr match { CALL strings_NewReader }
    gen { strings.NewReader($s); }
    bind { s as "s"; }
}
