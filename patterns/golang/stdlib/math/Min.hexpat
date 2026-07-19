arch x86_64;
platform linux, darwin, windows, freebsd;

pattern math_Min { name: "math.Min"; library: "go-stdlib"; instr match { CALL math_Min } gen { math.Min($x, $y); } bind { x as "x"; y as "y"; } }
