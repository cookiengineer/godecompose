arch x86_64;
platform linux, darwin, windows, freebsd;

pattern zip_Writer_Close { name: "zip.Writer.Close"; library: "go-stdlib"; instr match { CALL archive_zip___Writer__Close } gen { $w.Close(); } bind { w as "w"; } }
