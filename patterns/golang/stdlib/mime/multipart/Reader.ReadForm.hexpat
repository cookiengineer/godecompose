arch x86_64;
platform linux, darwin, windows, freebsd;

pattern multipart_Reader_ReadForm { name: "multipart.Reader.ReadForm"; library: "go-stdlib"; instr match { CALL mime_multipart___Reader__ReadForm } gen { $r.ReadForm($maxMemory); } bind { r as "r"; maxMemory as "maxMemory"; } }
