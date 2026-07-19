arch x86_64;
platform linux, darwin, windows, freebsd;

pattern path_filepath_Rel {
    name: "path/filepath.Rel"; library: "go-stdlib";
    instr match { CALL path_filepath_Rel }
    gen { filepath.Rel($basepath, $targpath); } bind { basepath as "basepath"; targpath as "targpath"; }
}
