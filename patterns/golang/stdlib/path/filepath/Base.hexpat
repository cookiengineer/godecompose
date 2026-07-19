arch x86_64;
platform linux, darwin, windows, freebsd;

pattern path_filepath_Base {
    name: "path/filepath.Base"; library: "go-stdlib";
    instr match { CALL path_filepath_Base }
    gen { filepath.Base($path); } bind { path as "path"; }
}
