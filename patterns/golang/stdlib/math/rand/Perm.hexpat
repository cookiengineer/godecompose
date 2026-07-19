arch x86_64;
platform linux, darwin, windows, freebsd;

pattern rand_Perm { name: "math/rand.Perm"; library: "go-stdlib"; instr match { CALL math_rand_Perm } gen { rand.Perm($n); } bind { n as "n"; } }
