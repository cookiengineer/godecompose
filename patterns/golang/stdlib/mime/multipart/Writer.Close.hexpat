arch x86_64;
platform linux, darwin, windows, freebsd;

pattern multipart_Writer_Close { name: "multipart.Writer.Close"; library: "go-stdlib"; instr match { CALL mime_multipart___Writer__Close } gen { $w.Close(); } bind { w as "w"; } }
