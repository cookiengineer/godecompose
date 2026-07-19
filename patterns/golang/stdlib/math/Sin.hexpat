arch x86_64;
platform linux, darwin, windows, freebsd;

pattern math_Sin { name: "math.Sin"; library: "go-stdlib"; instr match { CALL math_Sin } gen { math.Sin($x); } bind { x as "x"; } }
