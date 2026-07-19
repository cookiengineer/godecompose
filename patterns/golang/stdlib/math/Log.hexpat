arch x86_64;
platform linux, darwin, windows, freebsd;

pattern math_Log { name: "math.Log"; library: "go-stdlib"; instr match { CALL math_Log } gen { math.Log($x); } bind { x as "x"; } }
