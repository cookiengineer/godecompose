arch x86_64;
platform linux, darwin, windows, freebsd;

pattern zip_NewWriter { name: "archive/zip.NewWriter"; library: "go-stdlib"; instr match { CALL archive_zip_NewWriter } gen { zip.NewWriter($w); } bind { w as "w"; } }
