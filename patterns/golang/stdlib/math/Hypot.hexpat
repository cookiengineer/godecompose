arch x86_64;
platform linux, darwin, windows, freebsd;

pattern math_Hypot { name: "math.Hypot"; library: "go-stdlib"; instr match { CALL math_Hypot } gen { math.Hypot($p, $q); } bind { p as "p"; q as "q"; } }
