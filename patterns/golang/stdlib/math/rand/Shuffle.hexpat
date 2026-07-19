arch x86_64;
platform linux, darwin, windows, freebsd;

pattern rand_Shuffle { name: "math/rand.Shuffle"; library: "go-stdlib"; instr match { CALL math_rand_Shuffle } gen { rand.Shuffle($n, $swap); } bind { n as "n"; swap as "swap"; } }
