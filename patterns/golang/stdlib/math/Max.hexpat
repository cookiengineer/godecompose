arch x86_64;
platform linux, darwin, windows, freebsd;

pattern math_Max { name: "math.Max"; library: "go-stdlib"; instr match { CALL math_Max } gen { math.Max($x, $y); } bind { x as "x"; y as "y"; } }
