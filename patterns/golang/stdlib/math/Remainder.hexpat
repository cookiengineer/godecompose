arch x86_64;
platform linux, darwin, windows, freebsd;

pattern math_Remainder { name: "math.Remainder"; library: "go-stdlib"; instr match { CALL math_Remainder } gen { math.Remainder($x, $y); } bind { x as "x"; y as "y"; } }
