arch x86_64;
platform linux, darwin, windows, freebsd;

pattern gzip_NewWriter { name: "compress/gzip.NewWriter"; library: "go-stdlib"; instr match { CALL compress_gzip_NewWriter } gen { gzip.NewWriter($w); } bind { w as "w"; } }
