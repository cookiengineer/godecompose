arch x86_64;
platform linux, darwin, windows, freebsd;

pattern path_filepath_Split {
    name: "path/filepath.Split"; library: "go-stdlib";
    instr match { CALL path_filepath_Split }
    gen { filepath.Split($path); } bind { path as "path"; }
}
