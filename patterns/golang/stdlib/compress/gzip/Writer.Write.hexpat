arch x86_64;
platform linux, darwin, windows, freebsd;

pattern gzip_Writer_Write { name: "gzip.Writer.Write"; library: "go-stdlib"; instr match { CALL compress_gzip___Writer__Write } gen { $w.Write($p); } bind { w as "w"; p as "p"; } }
