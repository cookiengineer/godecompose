arch x86_64;
platform linux, darwin, windows, freebsd;

pattern mime_TypeByExtension { name: "mime.TypeByExtension"; library: "go-stdlib"; instr match { CALL mime_TypeByExtension } gen { mime.TypeByExtension($ext); } bind { ext as "ext"; } }
pattern mime_ExtensionsByType { name: "mime.ExtensionsByType"; library: "go-stdlib"; instr match { CALL mime_ExtensionsByType } gen { mime.ExtensionsByType($typ); } bind { typ as "typ"; } }
pattern mime_AddExtensionType { name: "mime.AddExtensionType"; library: "go-stdlib"; instr match { CALL mime_AddExtensionType } gen { mime.AddExtensionType($ext, $typ); } bind { ext as "ext"; typ as "typ"; } }
pattern mime_FormatMediaType { name: "mime.FormatMediaType"; library: "go-stdlib"; instr match { CALL mime_FormatMediaType } gen { mime.FormatMediaType($t, $param); } bind { t as "t"; param as "param"; } }
pattern mime_ParseMediaType { name: "mime.ParseMediaType"; library: "go-stdlib"; instr match { CALL mime_ParseMediaType } gen { mime.ParseMediaType($v); } bind { v as "v"; } }
