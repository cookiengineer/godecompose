arch x86_64;
platform linux, darwin, windows, freebsd;

pattern regexp_Match {
    name: "regexp.Match"; library: "go-stdlib";
    instr match { CALL regexp_Match }
    gen { regexp.Match($pattern, $b); } bind { pattern as "pattern"; b as "b"; }
}
