arch x86_64;
platform linux, darwin, windows, freebsd;

pattern tar_Writer_WriteHeader { name: "tar.Writer.WriteHeader"; library: "go-stdlib"; instr match { CALL archive_tar___Writer__WriteHeader } gen { $w.WriteHeader($hdr); } bind { w as "w"; hdr as "hdr"; } }
