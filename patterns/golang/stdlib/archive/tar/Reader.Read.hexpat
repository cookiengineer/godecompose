arch x86_64;
platform linux, darwin, windows, freebsd;

pattern tar_Reader_Read { name: "tar.Reader.Read"; library: "go-stdlib"; instr match { CALL archive_tar___Reader__Read } gen { $r.Read($b); } bind { r as "r"; b as "b"; } }
