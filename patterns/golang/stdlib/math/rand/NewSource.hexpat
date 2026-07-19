arch x86_64;
platform linux, darwin, windows, freebsd;

pattern rand_NewSource { name: "math/rand.NewSource"; library: "go-stdlib"; instr match { CALL math_rand_NewSource } gen { rand.NewSource($seed); } bind { seed as "seed"; } }
