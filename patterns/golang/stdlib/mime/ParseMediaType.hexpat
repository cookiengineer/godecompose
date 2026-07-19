arch x86_64;
platform linux, darwin, windows, freebsd;

pattern mime_ParseMediaType { name: "mime.ParseMediaType"; library: "go-stdlib"; instr match { CALL mime_ParseMediaType } gen { mime.ParseMediaType($v); } bind { v as "v"; } }
