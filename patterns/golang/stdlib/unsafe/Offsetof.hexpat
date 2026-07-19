arch x86_64;
platform linux, darwin, windows, freebsd;

pattern unsafe_Offsetof { name: "unsafe.Offsetof"; library: "go-stdlib"; instr match { CALL unsafe_Offsetof } gen { unsafe.Offsetof($x); } bind { x as "x"; } }
