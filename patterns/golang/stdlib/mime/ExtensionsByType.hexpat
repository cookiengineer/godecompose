arch x86_64;
platform linux, darwin, windows, freebsd;

pattern mime_ExtensionsByType { name: "mime.ExtensionsByType"; library: "go-stdlib"; instr match { CALL mime_ExtensionsByType } gen { mime.ExtensionsByType($typ); } bind { typ as "typ"; } }
