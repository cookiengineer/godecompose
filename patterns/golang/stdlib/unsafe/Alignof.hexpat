arch x86_64;
platform linux, darwin, windows, freebsd;

pattern unsafe_Alignof { name: "unsafe.Alignof"; library: "go-stdlib"; instr match { CALL unsafe_Alignof } gen { unsafe.Alignof($x); } bind { x as "x"; } }
