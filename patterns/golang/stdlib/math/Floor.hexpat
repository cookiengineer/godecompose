arch x86_64;
platform linux, darwin, windows, freebsd;

pattern math_Floor { name: "math.Floor"; library: "go-stdlib"; instr match { CALL math_Floor } gen { math.Floor($x); } bind { x as "x"; } }
