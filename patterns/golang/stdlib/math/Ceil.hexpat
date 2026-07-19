arch x86_64;
platform linux, darwin, windows, freebsd;

pattern math_Ceil { name: "math.Ceil"; library: "go-stdlib"; instr match { CALL math_Ceil } gen { math.Ceil($x); } bind { x as "x"; } }
