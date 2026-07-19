arch x86_64;
platform linux, darwin, windows, freebsd;

pattern multipart_NewWriter { name: "mime/multipart.NewWriter"; library: "go-stdlib"; instr match { CALL mime_multipart_NewWriter } gen { multipart.NewWriter($w); } bind { w as "w"; } }
