arch x86_64;
platform linux, darwin, windows, freebsd;

pattern zlib_NewWriterLevel { name: "compress/zlib.NewWriterLevel"; library: "go-stdlib"; instr match { CALL compress_zlib_NewWriterLevel } gen { zlib.NewWriterLevel($w, $level); } bind { w as "w"; level as "level"; } }
