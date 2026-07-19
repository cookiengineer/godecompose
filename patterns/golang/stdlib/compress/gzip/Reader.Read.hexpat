arch x86_64;
platform linux, darwin, windows, freebsd;

pattern gzip_Reader_Read { name: "gzip.Reader.Read"; library: "go-stdlib"; instr match { CALL compress_gzip___Reader__Read } gen { $r.Read($p); } bind { r as "r"; p as "p"; } }
