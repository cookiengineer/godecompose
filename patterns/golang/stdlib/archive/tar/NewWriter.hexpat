arch x86_64;
platform linux, darwin, windows, freebsd;

pattern tar_NewWriter { name: "archive/tar.NewWriter"; library: "go-stdlib"; instr match { CALL archive_tar_NewWriter } gen { tar.NewWriter($w); } bind { w as "w"; } }
