arch x86_64;
platform linux, darwin, windows, freebsd;

pattern mime_AddExtensionType { name: "mime.AddExtensionType"; library: "go-stdlib"; instr match { CALL mime_AddExtensionType } gen { mime.AddExtensionType($ext, $typ); } bind { ext as "ext"; typ as "typ"; } }
