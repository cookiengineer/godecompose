arch x86_64;
platform linux, darwin, windows, freebsd;

pattern zlib_NewReader { name: "compress/zlib.NewReader"; library: "go-stdlib"; instr match { CALL compress_zlib_NewReader } gen { zlib.NewReader($r); } bind { r as "r"; } }
