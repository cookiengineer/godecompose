arch x86_64;
platform linux, darwin, windows, freebsd;

pattern strings_ReplaceAll {
    name: "strings.ReplaceAll";
    library: "go-stdlib"; instr match { CALL strings_ReplaceAll }
    gen { strings.ReplaceAll($s, $old, $new); }
    bind { s as "s"; old as "old"; new as "new"; }
}
