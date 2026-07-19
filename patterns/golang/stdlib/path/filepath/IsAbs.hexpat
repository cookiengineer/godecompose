arch x86_64;
platform linux, darwin, windows, freebsd;

pattern path_filepath_IsAbs {
    name: "path/filepath.IsAbs"; library: "go-stdlib";
    instr match { CALL path_filepath_IsAbs }
    gen { filepath.IsAbs($path); } bind { path as "path"; }
}
