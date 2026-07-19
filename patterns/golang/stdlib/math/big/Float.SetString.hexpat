arch x86_64;
platform linux, darwin, windows, freebsd;

pattern big_Float_SetString { name: "math/big.Float.SetString"; library: "go-stdlib"; instr match { CALL math_big___Float__SetString } gen { $f.SetString($s); } bind { f as "f"; s as "s"; } }
