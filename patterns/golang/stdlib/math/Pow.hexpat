arch x86_64;
platform linux, darwin, windows, freebsd;

pattern math_Pow { name: "math.Pow"; library: "go-stdlib"; instr match { CALL math_Pow } gen { math.Pow($x, $y); } bind { x as "x"; y as "y"; } }
