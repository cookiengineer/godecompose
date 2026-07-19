arch x86_64;
platform linux, darwin, windows, freebsd;

pattern multipart_Reader_NextPart { name: "multipart.Reader.NextPart"; library: "go-stdlib"; instr match { CALL mime_multipart___Reader__NextPart } gen { $r.NextPart(); } bind { r as "r"; } }
