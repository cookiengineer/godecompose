arch x86_64;
platform linux, darwin, windows, freebsd;

pattern zip_File_Open { name: "zip.File.Open"; library: "go-stdlib"; instr match { CALL archive_zip___File__Open } gen { $f.Open(); } bind { f as "f"; } }
