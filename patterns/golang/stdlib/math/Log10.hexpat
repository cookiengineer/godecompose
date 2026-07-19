arch x86_64;
platform linux, darwin, windows, freebsd;

pattern math_Log10 { name: "math.Log10"; library: "go-stdlib"; instr match { CALL math_Log10 } gen { math.Log10($x); } bind { x as "x"; } }
