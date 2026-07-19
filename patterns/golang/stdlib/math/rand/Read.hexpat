arch x86_64;
platform linux, darwin, windows, freebsd;

pattern rand_Read { name: "math/rand.Read"; library: "go-stdlib"; instr match { CALL math_rand_Read } gen { rand.Read($p); } bind { p as "p"; } }
