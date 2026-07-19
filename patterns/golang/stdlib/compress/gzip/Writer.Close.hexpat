arch x86_64;
platform linux, darwin, windows, freebsd;

pattern gzip_Writer_Close { name: "gzip.Writer.Close"; library: "go-stdlib"; instr match { CALL compress_gzip___Writer__Close } gen { $w.Close(); } bind { w as "w"; } }
