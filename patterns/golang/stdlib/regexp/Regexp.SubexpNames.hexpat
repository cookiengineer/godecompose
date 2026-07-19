arch x86_64;
platform linux, darwin, windows, freebsd;

pattern regexp_Regexp_SubexpNames {
    name: "regexp.Regexp.SubexpNames"; library: "go-stdlib";
    instr match { CALL regexp___Regexp__SubexpNames }
    gen { $re.SubexpNames(); } bind { re as "re"; }
}
