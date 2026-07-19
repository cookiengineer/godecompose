arch x86_64;
platform linux, darwin, windows, freebsd;

pattern zip_Writer_Create { name: "zip.Writer.Create"; library: "go-stdlib"; instr match { CALL archive_zip___Writer__Create } gen { $w.Create($name); } bind { w as "w"; name as "name"; } }
