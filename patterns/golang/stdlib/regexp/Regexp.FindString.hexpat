arch x86_64;
platform linux, darwin, windows, freebsd;

pattern regexp_Regexp_FindString {
    name: "regexp.Regexp.FindString"; library: "go-stdlib";
    instr match { CALL regexp___Regexp__FindString }
    gen { $re.FindString($s); } bind { re as "re"; s as "s"; }
}
