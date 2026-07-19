arch x86_64;
platform linux, darwin, windows, freebsd;

pattern unsafe_SliceData { name: "unsafe.SliceData"; library: "go-stdlib"; instr match { CALL unsafe_SliceData } gen { unsafe.SliceData($s); } bind { s as "s"; } }
