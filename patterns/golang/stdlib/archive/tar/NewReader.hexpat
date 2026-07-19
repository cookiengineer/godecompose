arch x86_64;
platform linux, darwin, windows, freebsd;

pattern tar_NewReader { name: "archive/tar.NewReader"; library: "go-stdlib"; instr match { CALL archive_tar_NewReader } gen { tar.NewReader($r); } bind { r as "r"; } }
