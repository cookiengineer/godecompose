arch x86_64;
platform linux, darwin, windows, freebsd;

pattern regexp_CompilePOSIX {
    name: "regexp.CompilePOSIX"; library: "go-stdlib";
    instr match { CALL regexp_CompilePOSIX }
    gen { regexp.CompilePOSIX($expr); } bind { expr as "expr"; }
}
