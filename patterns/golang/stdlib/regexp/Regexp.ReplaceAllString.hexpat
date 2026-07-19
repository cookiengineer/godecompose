arch x86_64;
platform linux, darwin, windows, freebsd;

pattern regexp_Regexp_ReplaceAllString {
    name: "regexp.Regexp.ReplaceAllString"; library: "go-stdlib";
    instr match { CALL regexp___Regexp__ReplaceAllString }
    gen { $re.ReplaceAllString($src, $repl); } bind { re as "re"; src as "src"; repl as "repl"; }
}
