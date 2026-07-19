arch x86_64;
platform linux, darwin, windows, freebsd;

pattern unsafe_Slice { name: "unsafe.Slice"; library: "go-stdlib"; instr match { CALL unsafe_Slice } gen { unsafe.Slice($p, $len); } bind { p as "p"; len as "len"; } }
