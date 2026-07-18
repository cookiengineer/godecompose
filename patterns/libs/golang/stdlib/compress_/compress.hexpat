arch x86_64;
platform linux, darwin, windows, freebsd;

pattern gzip_NewReader { name: "compress/gzip.NewReader"; library: "go-stdlib"; instr match { CALL compress_gzip_NewReader } gen { gzip.NewReader($r); } bind { r as "r"; } }
pattern gzip_NewWriter { name: "compress/gzip.NewWriter"; library: "go-stdlib"; instr match { CALL compress_gzip_NewWriter } gen { gzip.NewWriter($w); } bind { w as "w"; } }
pattern gzip_NewWriterLevel { name: "compress/gzip.NewWriterLevel"; library: "go-stdlib"; instr match { CALL compress_gzip_NewWriterLevel } gen { gzip.NewWriterLevel($w, $level); } bind { w as "w"; level as "level"; } }
pattern gzip_Reader_Read { name: "gzip.Reader.Read"; library: "go-stdlib"; instr match { CALL compress_gzip___Reader__Read } gen { $r.Read($p); } bind { r as "r"; p as "p"; } }
pattern gzip_Reader_Close { name: "gzip.Reader.Close"; library: "go-stdlib"; instr match { CALL compress_gzip___Reader__Close } gen { $r.Close(); } bind { r as "r"; } }
pattern gzip_Writer_Write { name: "gzip.Writer.Write"; library: "go-stdlib"; instr match { CALL compress_gzip___Writer__Write } gen { $w.Write($p); } bind { w as "w"; p as "p"; } }
pattern gzip_Writer_Close { name: "gzip.Writer.Close"; library: "go-stdlib"; instr match { CALL compress_gzip___Writer__Close } gen { $w.Close(); } bind { w as "w"; } }
pattern gzip_Writer_Flush { name: "gzip.Writer.Flush"; library: "go-stdlib"; instr match { CALL compress_gzip___Writer__Flush } gen { $w.Flush(); } bind { w as "w"; } }
pattern gzip_Writer_Reset { name: "gzip.Writer.Reset"; library: "go-stdlib"; instr match { CALL compress_gzip___Writer__Reset } gen { $w.Reset($dst); } bind { w as "w"; dst as "dst"; } }

pattern zlib_NewReader { name: "compress/zlib.NewReader"; library: "go-stdlib"; instr match { CALL compress_zlib_NewReader } gen { zlib.NewReader($r); } bind { r as "r"; } }
pattern zlib_NewWriter { name: "compress/zlib.NewWriter"; library: "go-stdlib"; instr match { CALL compress_zlib_NewWriter } gen { zlib.NewWriter($w); } bind { w as "w"; } }
pattern zlib_NewWriterLevel { name: "compress/zlib.NewWriterLevel"; library: "go-stdlib"; instr match { CALL compress_zlib_NewWriterLevel } gen { zlib.NewWriterLevel($w, $level); } bind { w as "w"; level as "level"; } }

pattern flate_NewReader { name: "compress/flate.NewReader"; library: "go-stdlib"; instr match { CALL compress_flate_NewReader } gen { flate.NewReader($r); } bind { r as "r"; } }
pattern flate_NewWriter { name: "compress/flate.NewWriter"; library: "go-stdlib"; instr match { CALL compress_flate_NewWriter } gen { flate.NewWriter($w, $level); } bind { w as "w"; level as "level"; } }
pattern flate_NewWriterDict { name: "compress/flate.NewWriterDict"; library: "go-stdlib"; instr match { CALL compress_flate_NewWriterDict } gen { flate.NewWriterDict($w, $level, $dict); } bind { w as "w"; level as "level"; dict as "dict"; } }
