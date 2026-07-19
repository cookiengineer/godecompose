arch x86_64;
platform linux, darwin, windows, freebsd;

pattern unsafe_String { name: "unsafe.String"; library: "go-stdlib"; instr match { CALL unsafe_String } gen { unsafe.String($p, $len); } bind { p as "p"; len as "len"; } }
