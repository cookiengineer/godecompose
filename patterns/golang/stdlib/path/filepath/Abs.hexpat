arch x86_64;
platform linux, darwin, windows, freebsd;

pattern path_filepath_Abs {
    name: "path/filepath.Abs"; library: "go-stdlib";
    instr match { CALL path_filepath_Abs }
    gen { filepath.Abs($path); } bind { path as "path"; }
}
