arch x86_64;
platform linux, darwin, windows, freebsd;

pattern math_Tan { name: "math.Tan"; library: "go-stdlib"; instr match { CALL math_Tan } gen { math.Tan($x); } bind { x as "x"; } }
