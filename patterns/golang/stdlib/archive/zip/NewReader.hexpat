arch x86_64;
platform linux, darwin, windows, freebsd;

pattern zip_NewReader { name: "archive/zip.NewReader"; library: "go-stdlib"; instr match { CALL archive_zip_NewReader } gen { zip.NewReader($r, $size); } bind { r as "r"; size as "size"; } }
