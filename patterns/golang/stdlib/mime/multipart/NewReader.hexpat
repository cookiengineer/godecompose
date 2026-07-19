arch x86_64;
platform linux, darwin, windows, freebsd;

pattern multipart_NewReader { name: "mime/multipart.NewReader"; library: "go-stdlib"; instr match { CALL mime_multipart_NewReader } gen { multipart.NewReader($r, $boundary); } bind { r as "r"; boundary as "boundary"; } }
