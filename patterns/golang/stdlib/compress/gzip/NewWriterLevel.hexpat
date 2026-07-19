arch x86_64;
platform linux, darwin, windows, freebsd;

pattern gzip_NewWriterLevel { name: "compress/gzip.NewWriterLevel"; library: "go-stdlib"; instr match { CALL compress_gzip_NewWriterLevel } gen { gzip.NewWriterLevel($w, $level); } bind { w as "w"; level as "level"; } }
