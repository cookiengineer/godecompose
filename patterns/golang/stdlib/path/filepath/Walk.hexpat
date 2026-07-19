arch x86_64;
platform linux, darwin, windows, freebsd;

pattern path_filepath_Walk {
    name: "path/filepath.Walk"; library: "go-stdlib";
    instr match { CALL path_filepath_Walk }
    gen { filepath.Walk($root, $fn); } bind { root as "root"; fn as "fn"; }
}
