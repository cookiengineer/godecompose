arch x86_64;
platform linux, darwin, windows, freebsd;

pattern regexp_Regexp_Split {
    name: "regexp.Regexp.Split"; library: "go-stdlib";
    instr match { CALL regexp___Regexp__Split }
    gen { $re.Split($s, $n); } bind { re as "re"; s as "s"; n as "n"; }
}
