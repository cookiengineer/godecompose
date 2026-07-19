arch x86_64;
platform linux, darwin, windows, freebsd;

pattern gzip_Reader_Close { name: "gzip.Reader.Close"; library: "go-stdlib"; instr match { CALL compress_gzip___Reader__Close } gen { $r.Close(); } bind { r as "r"; } }
