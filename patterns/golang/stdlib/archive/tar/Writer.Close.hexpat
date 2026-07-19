arch x86_64;
platform linux, darwin, windows, freebsd;

pattern tar_Writer_Close { name: "tar.Writer.Close"; library: "go-stdlib"; instr match { CALL archive_tar___Writer__Close } gen { $w.Close(); } bind { w as "w"; } }
