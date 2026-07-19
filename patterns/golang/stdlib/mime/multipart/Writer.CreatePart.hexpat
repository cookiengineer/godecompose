arch x86_64;
platform linux, darwin, windows, freebsd;

pattern multipart_Writer_CreatePart { name: "multipart.Writer.CreatePart"; library: "go-stdlib"; instr match { CALL mime_multipart___Writer__CreatePart } gen { $w.CreatePart($header); } bind { w as "w"; header as "header"; } }
