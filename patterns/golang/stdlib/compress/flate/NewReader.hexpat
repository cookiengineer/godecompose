arch x86_64;
platform linux, darwin, windows, freebsd;

pattern flate_NewReader { name: "compress/flate.NewReader"; library: "go-stdlib"; instr match { CALL compress_flate_NewReader } gen { flate.NewReader($r); } bind { r as "r"; } }
