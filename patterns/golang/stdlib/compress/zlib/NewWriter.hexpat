arch x86_64;
platform linux, darwin, windows, freebsd;

pattern zlib_NewWriter { name: "compress/zlib.NewWriter"; library: "go-stdlib"; instr match { CALL compress_zlib_NewWriter } gen { zlib.NewWriter($w); } bind { w as "w"; } }
