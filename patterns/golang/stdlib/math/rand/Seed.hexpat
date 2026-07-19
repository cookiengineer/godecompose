arch x86_64;
platform linux, darwin, windows, freebsd;

pattern rand_Seed { name: "math/rand.Seed"; library: "go-stdlib"; instr match { CALL math_rand_Seed } gen { rand.Seed($seed); } bind { seed as "seed"; } }
