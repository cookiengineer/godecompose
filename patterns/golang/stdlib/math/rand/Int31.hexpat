arch x86_64;
platform linux, darwin, windows, freebsd;

pattern rand_Int31 { name: "math/rand.Int31"; library: "go-stdlib"; instr match { CALL math_rand_Int31 } gen { rand.Int31(); } }
