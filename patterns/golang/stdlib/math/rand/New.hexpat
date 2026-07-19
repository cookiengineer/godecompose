arch x86_64;
platform linux, darwin, windows, freebsd;

pattern rand_New { name: "math/rand.New"; library: "go-stdlib"; instr match { CALL math_rand_New } gen { rand.New($src); } bind { src as "src"; } }
