arch x86_64;
platform linux, darwin, windows, freebsd;

pattern regexp_Regexp_ReplaceAllStringFunc {
    name: "regexp.Regexp.ReplaceAllStringFunc"; library: "go-stdlib";
    instr match { CALL regexp___Regexp__ReplaceAllStringFunc }
    gen { $re.ReplaceAllStringFunc($src, $repl); } bind { re as "re"; src as "src"; repl as "repl"; }
}
