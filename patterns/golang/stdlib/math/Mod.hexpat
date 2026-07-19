arch x86_64;
platform linux, darwin, windows, freebsd;

pattern math_Mod { name: "math.Mod"; library: "go-stdlib"; instr match { CALL math_Mod } gen { math.Mod($x, $y); } bind { x as "x"; y as "y"; } }
