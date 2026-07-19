arch x86_64;
platform linux, darwin, windows, freebsd;

pattern unsafe_Sizeof { name: "unsafe.Sizeof"; library: "go-stdlib"; instr match { CALL unsafe_Sizeof } gen { unsafe.Sizeof($x); } bind { x as "x"; } }
