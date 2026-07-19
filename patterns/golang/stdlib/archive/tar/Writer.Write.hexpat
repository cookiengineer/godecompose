arch x86_64;
platform linux, darwin, windows, freebsd;

pattern tar_Writer_Write { name: "tar.Writer.Write"; library: "go-stdlib"; instr match { CALL archive_tar___Writer__Write } gen { $w.Write($b); } bind { w as "w"; b as "b"; } }
