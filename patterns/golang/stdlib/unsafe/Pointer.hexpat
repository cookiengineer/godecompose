arch x86_64;
platform linux, darwin, windows, freebsd;

pattern unsafe_Pointer { name: "unsafe.Pointer"; library: "go-stdlib"; instr match { CALL unsafe_Pointer } gen { unsafe.Pointer($p); } bind { p as "p"; } }
