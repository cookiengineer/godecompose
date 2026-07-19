arch x86_64;
platform linux, darwin, windows, freebsd;

pattern regexp_Regexp_FindAllString {
    name: "regexp.Regexp.FindAllString"; library: "go-stdlib";
    instr match { CALL regexp___Regexp__FindAllString }
    gen { $re.FindAllString($s, $n); } bind { re as "re"; s as "s"; n as "n"; }
}
