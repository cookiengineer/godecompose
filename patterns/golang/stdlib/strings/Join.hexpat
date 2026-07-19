arch x86_64;
platform linux, darwin, windows, freebsd;

pattern strings_Join {
    name: "strings.Join";
    library: "go-stdlib"; instr match { CALL strings_Join }
    gen { strings.Join($elems, $sep); }
    bind { elems as "elems"; sep as "sep"; }
}
