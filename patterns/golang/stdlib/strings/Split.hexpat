arch x86_64;
platform linux, darwin, windows, freebsd;

pattern strings_Split {
    name: "strings.Split";
    library: "go-stdlib"; instr match { CALL strings_Split }
    gen { strings.Split($s, $sep); }
    bind { s as "s"; sep as "sep"; }
}
