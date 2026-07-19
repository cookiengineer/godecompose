arch x86_64;
platform linux, darwin, windows, freebsd;

pattern rand_Int63 { name: "math/rand.Int63"; library: "go-stdlib"; instr match { CALL math_rand_Int63 } gen { rand.Int63(); } }
