arch x86_64;
platform linux, darwin, windows, freebsd;

pattern tar_Writer_Flush { name: "tar.Writer.Flush"; library: "go-stdlib"; instr match { CALL archive_tar___Writer__Flush } gen { $w.Flush(); } bind { w as "w"; } }
