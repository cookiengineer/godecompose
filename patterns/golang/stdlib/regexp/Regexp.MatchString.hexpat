arch x86_64;
platform linux, darwin, windows, freebsd;

pattern regexp_Regexp_MatchString {
    name: "regexp.Regexp.MatchString"; library: "go-stdlib";
    instr match { CALL regexp___Regexp__MatchString }
    gen { $re.MatchString($s); } bind { re as "re"; s as "s"; }
}
