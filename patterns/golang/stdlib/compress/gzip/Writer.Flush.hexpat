arch x86_64;
platform linux, darwin, windows, freebsd;

pattern gzip_Writer_Flush { name: "gzip.Writer.Flush"; library: "go-stdlib"; instr match { CALL compress_gzip___Writer__Flush } gen { $w.Flush(); } bind { w as "w"; } }
