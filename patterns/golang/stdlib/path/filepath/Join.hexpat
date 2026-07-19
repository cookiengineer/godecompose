arch x86_64;
platform linux, darwin, windows, freebsd;

pattern path_filepath_Join {
    name: "path/filepath.Join"; library: "go-stdlib";
    instr match { CALL path_filepath_Join }
    gen { filepath.Join($elem...); }
}
