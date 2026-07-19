arch x86_64;
platform linux, darwin, windows, freebsd;

pattern multipart_Writer_CreateFormFile { name: "multipart.Writer.CreateFormFile"; library: "go-stdlib"; instr match { CALL mime_multipart___Writer__CreateFormFile } gen { $w.CreateFormFile($fieldname, $filename); } bind { w as "w"; fieldname as "fieldname"; filename as "filename"; } }
