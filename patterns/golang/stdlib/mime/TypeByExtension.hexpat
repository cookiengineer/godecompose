arch x86_64;
platform linux, darwin, windows, freebsd;

pattern mime_TypeByExtension { name: "mime.TypeByExtension"; library: "go-stdlib"; instr match { CALL mime_TypeByExtension } gen { mime.TypeByExtension($ext); } bind { ext as "ext"; } }
