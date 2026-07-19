arch x86_64;
platform linux, darwin, windows, freebsd;

pattern path_filepath_Match {
    name: "path/filepath.Match"; library: "go-stdlib";
    instr match { CALL path_filepath_Match }
    gen { filepath.Match($pattern, $name); } bind { pattern as "pattern"; name as "name"; }
}
