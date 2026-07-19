arch x86_64;
platform linux, darwin, windows, freebsd;

pattern strings_TrimSuffix {
    name: "strings.TrimSuffix";
    library: "go-stdlib"; instr match { CALL strings_TrimSuffix }
    gen { strings.TrimSuffix($s, $suffix); }
    bind { s as "s"; suffix as "suffix"; }
}
