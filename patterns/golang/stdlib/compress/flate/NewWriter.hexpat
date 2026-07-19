arch x86_64;
platform linux, darwin, windows, freebsd;

pattern flate_NewWriter { name: "compress/flate.NewWriter"; library: "go-stdlib"; instr match { CALL compress_flate_NewWriter } gen { flate.NewWriter($w, $level); } bind { w as "w"; level as "level"; } }
