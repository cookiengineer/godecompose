arch x86_64;
platform linux, darwin, windows, freebsd;

pattern zip_OpenReader { name: "archive/zip.OpenReader"; library: "go-stdlib"; instr match { CALL archive_zip_OpenReader } gen { zip.OpenReader($name); } bind { name as "name"; } }
