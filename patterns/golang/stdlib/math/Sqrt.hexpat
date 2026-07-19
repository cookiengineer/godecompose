arch x86_64;
platform linux, darwin, windows, freebsd;

pattern math_Sqrt { name: "math.Sqrt"; library: "go-stdlib"; instr match { CALL math_Sqrt } gen { math.Sqrt($x); } bind { x as "x"; } }
