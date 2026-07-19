arch x86_64;
platform linux, darwin, windows, freebsd;

pattern regexp_QuoteMeta {
    name: "regexp.QuoteMeta"; library: "go-stdlib";
    instr match { CALL regexp_QuoteMeta }
    gen { regexp.QuoteMeta($s); } bind { s as "s"; }
}
