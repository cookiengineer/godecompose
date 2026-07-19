arch x86_64;
platform linux, darwin, windows, freebsd;

pattern path_filepath_Ext {
    name: "path/filepath.Ext"; library: "go-stdlib";
    instr match { CALL path_filepath_Ext }
    gen { filepath.Ext($path); } bind { path as "path"; }
}
