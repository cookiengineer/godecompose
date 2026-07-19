arch x86_64;
platform linux, darwin, windows, freebsd;

pattern mime_FormatMediaType { name: "mime.FormatMediaType"; library: "go-stdlib"; instr match { CALL mime_FormatMediaType } gen { mime.FormatMediaType($t, $param); } bind { t as "t"; param as "param"; } }
