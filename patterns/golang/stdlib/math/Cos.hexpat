arch x86_64;
platform linux, darwin, windows, freebsd;

pattern math_Cos { name: "math.Cos"; library: "go-stdlib"; instr match { CALL math_Cos } gen { math.Cos($x); } bind { x as "x"; } }
