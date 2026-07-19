arch x86_64;
platform linux, darwin, windows, freebsd;

pattern strings_HasSuffix {
    name: "strings.HasSuffix";
    library: "go-stdlib"; instr match { CALL strings_HasSuffix }
    gen { strings.HasSuffix($s, $suffix); }
    bind { s as "s"; suffix as "suffix"; }
}
