arch x86_64;
platform linux, darwin, windows, freebsd;

pattern flate_NewWriterDict { name: "compress/flate.NewWriterDict"; library: "go-stdlib"; instr match { CALL compress_flate_NewWriterDict } gen { flate.NewWriterDict($w, $level, $dict); } bind { w as "w"; level as "level"; dict as "dict"; } }
