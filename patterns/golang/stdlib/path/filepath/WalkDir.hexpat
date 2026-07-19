arch x86_64;
platform linux, darwin, windows, freebsd;

pattern path_filepath_WalkDir {
    name: "path/filepath.WalkDir"; library: "go-stdlib";
    instr match { CALL path_filepath_WalkDir }
    gen { filepath.WalkDir($root, $fn); } bind { root as "root"; fn as "fn"; }
}
