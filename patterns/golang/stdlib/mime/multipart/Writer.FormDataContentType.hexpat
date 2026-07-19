arch x86_64;
platform linux, darwin, windows, freebsd;

pattern multipart_Writer_FormDataContentType { name: "multipart.Writer.FormDataContentType"; library: "go-stdlib"; instr match { CALL mime_multipart___Writer__FormDataContentType } gen { $w.FormDataContentType(); } bind { w as "w"; } }
