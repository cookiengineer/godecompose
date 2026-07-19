arch x86_64;
platform linux, darwin, windows, freebsd;

pattern math_Abs { name: "math.Abs"; library: "go-stdlib"; instr match { CALL math_Abs } gen { math.Abs($x); } bind { x as "x"; } }
