arch x86_64;
platform linux, darwin, windows, freebsd;

pattern gzip_NewReader { name: "compress/gzip.NewReader"; library: "go-stdlib"; instr match { CALL compress_gzip_NewReader } gen { gzip.NewReader($r); } bind { r as "r"; } }
