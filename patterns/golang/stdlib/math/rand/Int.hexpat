arch x86_64;
platform linux, darwin, windows, freebsd;

pattern rand_Int { name: "math/rand.Int"; library: "go-stdlib"; instr match { CALL math_rand_Int } gen { rand.Int(); } }
