arch x86_64;
platform linux, darwin, windows, freebsd;

pattern tar_Reader_Next { name: "tar.Reader.Next"; library: "go-stdlib"; instr match { CALL archive_tar___Reader__Next } gen { $r.Next(); } bind { r as "r"; } }
