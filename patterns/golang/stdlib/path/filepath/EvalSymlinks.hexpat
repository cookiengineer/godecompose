arch x86_64;
platform linux, darwin, windows, freebsd;

pattern path_filepath_EvalSymlinks {
    name: "path/filepath.EvalSymlinks"; library: "go-stdlib";
    instr match { CALL path_filepath_EvalSymlinks }
    gen { filepath.EvalSymlinks($path); } bind { path as "path"; }
}
