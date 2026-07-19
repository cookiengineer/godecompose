arch x86_64;
platform linux, darwin, windows, freebsd;

pattern rand_Float64 { name: "math/rand.Float64"; library: "go-stdlib"; instr match { CALL math_rand_Float64 } gen { rand.Float64(); } }
