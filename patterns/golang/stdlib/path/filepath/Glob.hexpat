arch x86_64;
platform linux, darwin, windows, freebsd;

pattern path_filepath_Glob {
    name: "path/filepath.Glob"; library: "go-stdlib";
    instr match { CALL path_filepath_Glob }
    gen { filepath.Glob($pattern); } bind { pattern as "pattern"; }
}
