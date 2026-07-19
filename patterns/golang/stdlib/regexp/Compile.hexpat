arch x86_64;
platform linux, darwin, windows, freebsd;

pattern regexp_Compile {
    name: "regexp.Compile"; library: "go-stdlib";
    instr match { CALL regexp_Compile }
    gen { regexp.Compile($expr); } bind { expr as "expr"; }
}
