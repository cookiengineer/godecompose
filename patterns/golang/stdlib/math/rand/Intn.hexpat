arch x86_64;
platform linux, darwin, windows, freebsd;

pattern rand_Intn { name: "math/rand.Intn"; library: "go-stdlib"; instr match { CALL math_rand_Intn } gen { rand.Intn($n); } bind { n as "n"; } }
