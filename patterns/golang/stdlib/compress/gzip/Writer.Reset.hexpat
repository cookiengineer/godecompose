arch x86_64;
platform linux, darwin, windows, freebsd;

pattern gzip_Writer_Reset { name: "gzip.Writer.Reset"; library: "go-stdlib"; instr match { CALL compress_gzip___Writer__Reset } gen { $w.Reset($dst); } bind { w as "w"; dst as "dst"; } }
