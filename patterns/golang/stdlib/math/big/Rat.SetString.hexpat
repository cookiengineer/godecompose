arch x86_64;
platform linux, darwin, windows, freebsd;

pattern big_Rat_SetString { name: "math/big.Rat.SetString"; library: "go-stdlib"; instr match { CALL math_big___Rat__SetString } gen { $r.SetString($s); } bind { r as "r"; s as "s"; } }
