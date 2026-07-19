arch x86_64;
platform linux, darwin, windows, freebsd;

pattern path_filepath_Clean {
    name: "path/filepath.Clean"; library: "go-stdlib";
    instr match { CALL path_filepath_Clean }
    gen { filepath.Clean($path); } bind { path as "path"; }
}
