arch x86_64;
platform linux, darwin, windows, freebsd;

pattern math_Exp { name: "math.Exp"; library: "go-stdlib"; instr match { CALL math_Exp } gen { math.Exp($x); } bind { x as "x"; } }
