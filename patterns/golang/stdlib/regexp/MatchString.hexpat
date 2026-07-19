arch x86_64;
platform linux, darwin, windows, freebsd;

pattern regexp_MatchString {
    name: "regexp.MatchString"; library: "go-stdlib";
    instr match { CALL regexp_MatchString }
    gen { regexp.MatchString($pattern, $s); } bind { pattern as "pattern"; s as "s"; }
}
