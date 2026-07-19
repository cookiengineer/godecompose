arch x86_64;
platform linux, darwin, windows, freebsd;

pattern math_Round { name: "math.Round"; library: "go-stdlib"; instr match { CALL math_Round } gen { math.Round($x); } bind { x as "x"; } }
