arch x86_64;
platform linux, darwin, windows, freebsd;

pattern math_Log2 { name: "math.Log2"; library: "go-stdlib"; instr match { CALL math_Log2 } gen { math.Log2($x); } bind { x as "x"; } }
