arch x86_64;
platform linux, darwin, windows, freebsd;

pattern regexp_MustCompile {
    name: "regexp.MustCompile"; library: "go-stdlib";
    instr match { CALL regexp_MustCompile }
    gen { regexp.MustCompile($expr); } bind { expr as "expr"; }
}
