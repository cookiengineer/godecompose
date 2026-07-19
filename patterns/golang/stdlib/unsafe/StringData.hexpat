arch x86_64;
platform linux, darwin, windows, freebsd;

pattern unsafe_StringData { name: "unsafe.StringData"; library: "go-stdlib"; instr match { CALL unsafe_StringData } gen { unsafe.StringData($s); } bind { s as "s"; } }
