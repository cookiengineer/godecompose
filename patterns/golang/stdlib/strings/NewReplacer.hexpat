arch x86_64;
platform linux, darwin, windows, freebsd;

pattern strings_NewReplacer {
    name: "strings.NewReplacer";
    library: "go-stdlib"; instr match { CALL strings_NewReplacer }
    gen { strings.NewReplacer($oldnew...); }
    bind { oldnew as "oldnew"; }
}
