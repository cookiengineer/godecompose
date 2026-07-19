arch x86_64;
platform linux, darwin, windows, freebsd;

pattern path_filepath_Dir {
    name: "path/filepath.Dir"; library: "go-stdlib";
    instr match { CALL path_filepath_Dir }
    gen { filepath.Dir($path); } bind { path as "path"; }
}
