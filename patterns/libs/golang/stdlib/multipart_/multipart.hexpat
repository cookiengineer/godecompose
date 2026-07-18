arch x86_64;
platform linux, darwin, windows, freebsd;

pattern multipart_NewReader { name: "mime/multipart.NewReader"; library: "go-stdlib"; instr match { CALL mime_multipart_NewReader } gen { multipart.NewReader($r, $boundary); } bind { r as "r"; boundary as "boundary"; } }
pattern multipart_NewWriter { name: "mime/multipart.NewWriter"; library: "go-stdlib"; instr match { CALL mime_multipart_NewWriter } gen { multipart.NewWriter($w); } bind { w as "w"; } }
pattern multipart_Reader_NextPart { name: "multipart.Reader.NextPart"; library: "go-stdlib"; instr match { CALL mime_multipart___Reader__NextPart } gen { $r.NextPart(); } bind { r as "r"; } }
pattern multipart_Reader_ReadForm { name: "multipart.Reader.ReadForm"; library: "go-stdlib"; instr match { CALL mime_multipart___Reader__ReadForm } gen { $r.ReadForm($maxMemory); } bind { r as "r"; maxMemory as "maxMemory"; } }
pattern multipart_Writer_CreateFormFile { name: "multipart.Writer.CreateFormFile"; library: "go-stdlib"; instr match { CALL mime_multipart___Writer__CreateFormFile } gen { $w.CreateFormFile($fieldname, $filename); } bind { w as "w"; fieldname as "fieldname"; filename as "filename"; } }
pattern multipart_Writer_CreatePart { name: "multipart.Writer.CreatePart"; library: "go-stdlib"; instr match { CALL mime_multipart___Writer__CreatePart } gen { $w.CreatePart($header); } bind { w as "w"; header as "header"; } }
pattern multipart_Writer_FormDataContentType { name: "multipart.Writer.FormDataContentType"; library: "go-stdlib"; instr match { CALL mime_multipart___Writer__FormDataContentType } gen { $w.FormDataContentType(); } bind { w as "w"; } }
pattern multipart_Writer_Close { name: "multipart.Writer.Close"; library: "go-stdlib"; instr match { CALL mime_multipart___Writer__Close } gen { $w.Close(); } bind { w as "w"; } }
