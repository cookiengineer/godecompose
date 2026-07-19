arch x86_64;
platform linux, darwin, windows, freebsd;

pattern zip_Writer_CreateHeader { name: "zip.Writer.CreateHeader"; library: "go-stdlib"; instr match { CALL archive_zip___Writer__CreateHeader } gen { $w.CreateHeader($fh); } bind { w as "w"; fh as "fh"; } }
